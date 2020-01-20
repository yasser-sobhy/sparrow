package core

import (
	"bytes"
	"io"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/evio"
)

type Conn struct {
	Evio     evio.Conn
	upgraded bool
	ctx      interface{}
}

func (c *Conn) Context() interface{}       { return c.ctx }
func (c *Conn) SetContext(ctx interface{}) { c.ctx = ctx }

func (c *Conn) Write(message []byte) {

}

func (c *Conn) Close() {

}

func (c *Conn) Ping() {

}

type ReadWriteBuffer struct {
	r *bytes.Buffer
	w *bytes.Buffer
}

func (rw ReadWriteBuffer) Read(p []byte) (n int, err error) {
	return rw.r.Read(p)
}

func (rw ReadWriteBuffer) Write(p []byte) (n int, err error) {
	return rw.w.Write(p)
}

type NanoWebsocket struct {
	OnOpen    func(c *Conn, handshake *ws.Handshake)
	OnClose   func(c *Conn, err error) (action evio.Action)
	OnMessage func(c *Conn, mesaage wsutil.Message) (out []byte, action evio.Action)
	OnPing    func(c *Conn)
	OnPong    func(c *Conn)
}

func (nano *NanoWebsocket) Serve(addr ...string) error {

	var events evio.Events
	events.Opened = func(conn evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		c := Conn{Evio: conn, upgraded: false}
		conn.SetContext(c)
		return
	}

	events.Closed = func(conn evio.Conn, err error) (action evio.Action) {
		c, ok := conn.Context().(Conn)
		if ok {
			return nano.OnClose(&c, err)
		}

		return
	}

	events.Data = func(conn evio.Conn, in []byte) (out []byte, action evio.Action) {

		c, ok := conn.Context().(Conn)

		if !ok {
			logrus.Error("couldn't assert connection", conn.RemoteAddr().String())
			return
		}

		if c.upgraded {
			return nano.processData(&c, in)
		}

		rwb := ReadWriteBuffer{
			r: bytes.NewBuffer(in),
			w: bytes.NewBuffer(out),
		}

		handshake, err := ws.Upgrade(rwb)
		if err != nil {
			action = evio.Close
			logrus.Error(err)
			return
		}

		out = rwb.w.Bytes()
		nano.OnOpen(&c, &handshake)
		return
	}

	return evio.Serve(events, addr...)
}

func (nano *NanoWebsocket) processData(conn *Conn, in []byte) (out []byte, action evio.Action) {

	rwb := ReadWriteBuffer{
		r: bytes.NewBuffer(in),
		w: bytes.NewBuffer(out),
	}

	controlHandler := wsutil.ControlFrameHandler(rwb, ws.StateServerSide)

	rd := wsutil.Reader{
		Source:         rwb.r,
		State:          ws.StateServerSide,
		CheckUTF8:      true,
		OnIntermediate: controlHandler,
	}

	h, err := rd.NextFrame()
	if err != nil {
		action = evio.Close
		return
	}

	var p []byte
	if h.Fin {
		// No more frames will be read. Use fixed sized buffer to read payload.
		p = make([]byte, h.Length)
		// It is not possible to receive io.EOF here because Reader does not
		// return EOF if frame payload was successfully fetched.
		// Thus we consistent here with io.Reader behavior.
		_, err = io.ReadFull(&rd, p)
	} else {
		// Frame is fragmented, thus use ioutil.ReadAll behavior.
		var buf bytes.Buffer
		_, err = buf.ReadFrom(&rd)
		p = buf.Bytes()
	}

	if err != nil {
		action = evio.Close
		return
	}

	return nano.OnMessage(conn, wsutil.Message{h.OpCode, p})
}