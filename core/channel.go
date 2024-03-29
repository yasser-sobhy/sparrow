package core

type Channel struct {
	ID string
	// display name
	Name string

	subscribers   map[*Conn]struct{}
	sharedMessage []byte
}

func (channel *Channel) Join(user *Conn) {
	channel.subscribers[user] = struct{}{}
}

func (channel *Channel) Leave(user *Conn) {
	delete(channel.subscribers, user)
}

func (channel *Channel) Send(data []byte) {
	if len(channel.subscribers) > 0 {
		channel.sharedMessage = append(channel.sharedMessage, data...)
		defaultChannelLoop.Add(channel)
	}
}

type channelLoop struct {
	channels []*Channel
}

var defaultChannelLoop channelLoop

func (channelLoop *channelLoop) Add(channel *Channel) {
	channelLoop.channels = append(channelLoop.channels, channel)
}

func (channelLoop *channelLoop) run() {
	go func() {

		if len(channelLoop.channels) == 0 {
			return
		}

		for _, channel := range channelLoop.channels {
			go func() {
				for ws, _ := range channel.subscribers {
					ws.Write(channel.sharedMessage)
				}
			}()
			channel.sharedMessage = nil
		}
		channelLoop.channels = nil
	}()

}
