package crmon

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crmon/pkg/log"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

type app struct {
	opts        Options
	logger      log.ZeroLogger
	subscribers []Subscriber
}

func (a *app) Run() (err error) {
	ctx := context.TODO()
	client, err := pubsub.NewClient(ctx, a.opts.ProjectID)
	if err != nil {
		return errors.Wrap(err, "create pubsub client")
	}

	topic, err := a.ensureTopic(ctx, client, a.opts.Topic)
	if err != nil {
		return errors.Wrap(err, "ensure topic")
	}

	sub, err := a.ensureSubscription(ctx, client, topic, a.opts.Subscription)
	if err != nil {
		return errors.Wrap(err, "ensure subscription")
	}

	return a.monitor(ctx, sub)
}

func (a *app) ensureTopic(ctx context.Context, client *pubsub.Client, topicID string) (topic *pubsub.Topic, err error) {
	topic = client.Topic(topicID)

	// check if topic exists
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "check exists")
	}
	if exists {
		a.logger.Info().Str("topic", topicID).Msg("topic exists")
		return
	}

	// create new topic
	a.logger.Info().
		Str("topic", topicID).
		Msg("creating new topic")
	topic, err = client.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, errors.Wrap(err, "create topic")
	}
	return
}

func (a *app) ensureSubscription(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, subscriptionID string) (sub *pubsub.Subscription, err error) {
	sub = client.Subscription(subscriptionID)

	// check if subscription exists
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "check exists")
	}
	if exists {
		a.logger.Info().Str("subscription", subscriptionID).Msg("subscription exists")
		return
	}

	// create new subscription
	a.logger.Info().
		Str("topic", topic.ID()).
		Str("subscription", subscriptionID).
		Msg("create new subscription")
	sub, err = client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create subscription")
	}
	return
}

func (a *app) monitor(ctx context.Context, sub *pubsub.Subscription) (err error) {
	for _, s := range a.subscribers {
		if err = s.Init(); err != nil {
			a.logger.Warn().Msg("cannot init " + s.Name() + " subscriber")
		}
	}
	defer func() {
		for _, s := range a.subscribers {
			if err = s.Cleanup(); err != nil {
				a.logger.Warn().Msg("cannot init " + s.Name() + " subscriber")
			}
		}
	}()

	a.logger.Info().
		Str("subscription", sub.ID()).
		Msg("listening for new image updates")

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		a.logger.Debug().Str("id", msg.ID).Msg("received msg")

		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			a.logger.Warn().
				Err(err).
				Str("id", msg.ID).
				Str("data", string(msg.Data)).
				Msg("unmarshal failed")
			msg.Ack()
			return
		}

		for _, s := range a.subscribers {
			if err := s.OnReceive(event); err != nil {
				a.logger.Warn().Err(err).Msg(s.Name() + " failed to receive new event")
			}
		}

		// temporary disable for testing
		if os.Getenv("CRMON_TESTING") != "" {
			return
		}
		msg.Ack()
	})
	return err
}

func NewApp(opts Options) App {
	return &app{
		opts:        opts,
		logger:      log.NewZeroLogger(),
		subscribers: opts.Subscribers,
	}
}
