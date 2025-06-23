package server

type Broker struct {
	subscribers map[string]map[*Client]bool
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]map[*Client]bool),
	}
}

func (b *Broker) Subscribe(topic string, client *Client) {
	if b.subscribers[topic] == nil {
		b.subscribers[topic] = make(map[*Client]bool)
	}
	b.subscribers[topic][client] = true
}

func (b *Broker) Unsubscribe(topic string, client *Client) {
	if subs, ok := b.subscribers[topic]; ok {
		delete(subs, client)
	}
}

func (b *Broker) UnsubscribeAll(client *Client) {
	for topic, subs := range b.subscribers {
		if subs[client] {
			delete(subs, client)
			if len(subs) == 0 {
				delete(b.subscribers, topic)
			}
		}
	}
}

func (b *Broker) Publish(topic string, msg string) {
	for client := range b.subscribers[topic] {
		client.SendMessage(topic, msg)
	}
}
