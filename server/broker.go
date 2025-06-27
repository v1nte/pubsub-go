package server

import "log"

type subscriptionRequest struct {
	client *Client
	topic  string
}

type unsubscriptionRequest struct {
	client *Client
	topic  string
}

type publishRequest struct {
	author  string
	topic   string
	message string
}

type Broker struct {
	subscribers map[string]map[*Client]bool

	subscribeChan   chan subscriptionRequest
	unsubscribeChan chan unsubscriptionRequest
	publishChan     chan publishRequest
	unsubscribeAll  chan *Client
}

func NewBroker() *Broker {
	b := &Broker{
		subscribers:     make(map[string]map[*Client]bool),
		subscribeChan:   make(chan subscriptionRequest),
		unsubscribeChan: make(chan unsubscriptionRequest),
		publishChan:     make(chan publishRequest),
		unsubscribeAll:  make(chan *Client),
	}

	go b.run()
	return b
}

func (b *Broker) run() {
	for {
		select {
		case sub := <-b.subscribeChan:
			if b.subscribers[sub.topic] == nil {
				b.subscribers[sub.topic] = make(map[*Client]bool)
			}
			b.subscribers[sub.topic][sub.client] = true

		case unsub := <-b.unsubscribeChan:
			if subs, ok := b.subscribers[unsub.topic]; ok {
				delete(subs, unsub.client)
				if len(subs) == 0 {
					delete(b.subscribers, unsub.topic)
				}

			}

		case pub := <-b.publishChan:
			for client := range b.subscribers[pub.topic] {
				select {
				case client.send <- OutgoingMessage{
					Author:  pub.author,
					Topic:   pub.topic,
					Message: pub.message,
				}:
				default:
					log.Println("Channel full")
				}
			}

		case client := <-b.unsubscribeAll:
			for topic, sub := range b.subscribers {
				if sub[client] {
					delete(sub, client)
					if len(sub) == 0 {
						delete(b.subscribers, topic)
					}
				}
			}
		}
	}
}
