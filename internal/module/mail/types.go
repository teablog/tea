package mail

type Message struct {
	To          []string
	Title       string
	Body        string
	ContentType string
}

const (
	ContentTypeText = "text/plain"
	ContentTypeHtml = "text/html"
)

func NewMessage() *Message {
	return &Message{
		ContentType: ContentTypeHtml,
	}
}

func (m *Message) SetTo(to ...string) *Message {
	m.To = to
	return m
}

func (m *Message) SetTitle(t string) *Message {
	m.Title = t
	return m
}

func (m *Message) SetBody(t string) *Message {
	m.Body = t
	return m
}

func (m *Message) SetContentType(ct string) *Message {
	m.ContentType = ct
	return m
}

