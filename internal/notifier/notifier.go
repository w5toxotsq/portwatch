package notifier

// Notifier defines the interface for sending port change notifications.
type Notifier interface {
	Notify(event Event) error
}

// EventType represents the type of port change event.
type EventType string

const (
	EventOpened EventType = "opened"
	EventClosed EventType = "closed"
)

// Event represents a single port change notification.
type Event struct {
	Type     EventType
	Protocol string
	Port     int
	Host     string
}

// Multi dispatches notifications to multiple Notifier implementations.
type Multi struct {
	notifiers []Notifier
}

// NewMulti creates a Multi notifier that fans out to all provided notifiers.
func NewMulti(notifiers ...Notifier) *Multi {
	return &Multi{notifiers: notifiers}
}

// Notify sends the event to all registered notifiers, collecting errors.
func (m *Multi) Notify(event Event) error {
	var firstErr error
	for _, n := range m.notifiers {
		if err := n.Notify(event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
