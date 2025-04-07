package notifications

type Channel string

const (
	ChannelEmail Channel = "email"
)

const (
	TypeActivationEmail = "activation"
	TypePasswordChanged = "password_changed"
)

type Message struct {
	Channel Channel `json:"channel"`
	Type    string  `json:"type"`
	Content any     `json:"content"`
}

func NewMessage(channel Channel, _type string, content any) *Message {
	return &Message{Channel: channel, Type: _type, Content: content}
}
