package broker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack/v4"
	"go.uber.org/fx"
)

type broker struct {
	subscription stan.Subscription
	client       stan.Conn
	config       *viper.Viper
}

func (b *broker) callback(msg *stan.Msg) {
	dec := msgpack.NewDecoder(bytes.NewBuffer(msg.Data))
	logger.Debug().Msg(fmt.Sprintf("handling msg from %s:%d...", msg.Subject, msg.Sequence))
	msgType, err := dec.Query("*.message.type")
	if err != nil {
		return
	}

	switch msgType[0].(string) {
	case "text":
		_ = b.handleTextMsg(msg)
	}

	_ = msg.Ack()
}

func (b *broker) handleTextMsg(msg *stan.Msg) error {
	logger.Debug().Msg("handling text msg...")
	channel := b.config.GetString("text_channel")
	if err := b.client.Publish(channel, msg.Data); err != nil {
		return err
	}
	return nil
}

func (b *broker) run() error {
	channel := b.config.GetString("msg_channel")
	var err error
	b.subscription, err = b.client.Subscribe(channel, b.callback)
	return err
}

func (b *broker) stop() error {
	if err := b.subscription.Close(); err != nil {
		return err
	}
	if err := b.client.Close(); err != nil {
		return err
	}
	return nil
}

func newBroker(lc fx.Lifecycle, config *viper.Viper, client stan.Conn) (*broker, error) {
	broker := broker{
		client: client,
		config: config,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return broker.run()
		},
		OnStop: func(ctx context.Context) error {
			return broker.stop()
		},
	})

	return &broker, nil
}
