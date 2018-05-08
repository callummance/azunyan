package manager

import broadcast "github.com/dustin/go-broadcast"

const broadcastBufLen int = 10

type BroadcastData struct {
	Name    string
	Content interface{}
}

func InitBroadcaster() broadcast.Broadcaster {
	return broadcast.NewBroadcaster(broadcastBufLen)
}

func (m *KaraokeManager) SubscribeToChanges() chan interface{} {
	newListener := make(chan interface{}, broadcastBufLen)
	m.UpdateSubscribers.Register(newListener)

	return newListener
}

func (m *KaraokeManager) Unsubscribe(target chan interface{}) {
	m.UpdateSubscribers.Unregister(target)
}

func (m *KaraokeManager) SendBroadcast(bc BroadcastData) {
	m.UpdateSubscribers.Submit(bc)
}
