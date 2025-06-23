package server

type Broker struct {
	suscribers map[string]map[*Client]bool
}

func NewBroker() *Broker {
	return &Broker{
		suscribers: make(map[string]map[*Client]bool),
	}
}

func (b *Broker) Suscribe(topic string, client *Client) {
	if b.suscribers[topic] == nil {
		b.suscribers[topic] = make(map[*Client]bool)
	}
	b.suscribers[topic][client] = true
}

func (b *Broker) Unsuscribe(topic string, client *Client) {
	if subs, ok := b.suscribers[topic]; ok {
		delete(subs, client)
	}
}

func (b *Broker) Publish(topic string, msg string) {
	for client := range b.suscribers[topic] {
		client.SendMessage(topic, msg)
	}
}
