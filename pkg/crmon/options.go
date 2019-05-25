package crmon

type Options struct {
	ProjectID    string
	Topic        string
	Subscription string
	Subscribers  []Subscriber
}
