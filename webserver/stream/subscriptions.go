package stream

import (
	"github.com/dustin/go-broadcast"
)

const broadcastBufLen int = 10
var subscriptions broadcast.Broadcaster

type BroadcastData struct {
	Name string
	Content interface{}
}




func InitBroadcaster() {
	subscriptions = broadcast.NewBroadcaster(broadcastBufLen)
}

func SubscribeToChanges() chan interface{}{
	newListener := make(chan interface{}, broadcastBufLen)
	subscriptions.Register(newListener)

	return newListener
}

func Unsubscribe(target chan interface{}){
	subscriptions.Unregister(target)
}

func SendBroadcast(bc BroadcastData) {
	subscriptions.Submit(bc)
}