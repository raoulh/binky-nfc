package acr122u

// Handler is the interface that handles each card when present in the field.
type Handler interface {
	ServeCard(Card)
	CardRemoved()
}

// HandlerFunc is the function signature for handling a card
type HandlerFunc func(Card)

// ServeCard makes HandlerFunc implement the Handler interface
func (hf HandlerFunc) ServeCard(c Card) {
	hf(c)
}

func (hf HandlerFunc) CardRemoved() {
}