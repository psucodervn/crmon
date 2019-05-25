package crmon

type App interface {
	Run() error
}

type Subscriber interface {
	Name() string
	Init() error
	Cleanup() error
	OnReceive(event Event) error
}

type Event struct {
	Action string `json:"action"`
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
}

const (
	ActionInsert = "INSERT"
	ActionDelete = "DELETE"
)
