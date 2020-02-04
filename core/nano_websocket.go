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
	evio.Conn

	upgraded bool
	out      []byte
}

func (c *Conn) Write(message []byte) {
	c.out = message
	c.Wake()
}

func (c *Conn) WriteMessage(message []byte) {
	m, _ := ws.CompileFrame(ws.NewTextFrame(message))
	c.Write(m)
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
		c := Conn{conn, false, nil}
		conn.SetContext(&c)
		return
	}

	events.Closed = func(conn evio.Conn, err error) (action evio.Action) {
		c, ok := conn.Context().(*Conn)
		if ok {
			return nano.OnClose(c, err)
		}

		return
	}

	events.Data = func(conn evio.Conn, in []byte) (out []byte, action evio.Action) {

		c, ok := conn.Context().(*Conn)

		if !ok {
			logrus.Error("couldn't assert connection", conn.RemoteAddr().String())
			return
		}

		if c.upgraded {
			// sending data
			if in == nil {
				out = c.out
				return
			}

			return nano.processData(c, in)
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
		c.upgraded = true
		nano.OnOpen(c, &handshake)
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
		logrus.Error("closig conn. failed to read NextFrame ")
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
		logrus.Error("closig conn. failed to read NextFrame ")
		action = evio.Close
		return
	}

	return nano.OnMessage(conn, wsutil.Message{h.OpCode, p})
}
