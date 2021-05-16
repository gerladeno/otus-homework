package rmq

import (
	"context"
	"fmt"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Client struct {
	log  *logrus.Logger
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
	busy chan struct{}
}

const retry = 5

func GetRMQConnectionAndDeclare(log *logrus.Logger, dsn string, ttl int64) (*Client, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	var args amqp.Table
	if ttl != 0 {
		args = amqp.Table{"x-message-ttl": ttl}
	}
	topic, err := ch.QueueDeclare("notifications", false, false, false, false, args)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	return &Client{
		log:  log,
		conn: conn,
		ch:   ch,
		q:    topic,
		busy: make(chan struct{}, 1),
	}, nil
}

func (c *Client) Close() error {
	if err := c.ch.Close(); err != nil {
		c.log.Warn("err closing rmq channel: ", err)
	}
	<-c.busy
	close(c.busy)
	return c.conn.Close()
}

func (c *Client) Notify(events []common.Event) {
	for _, event := range events {
		msg, err := event.Notification().Encode()
		if err != nil {
			c.log.Warn("failed to encode msg: ", event.Notification().String())
			continue
		}
		for i := 0; i < retry; i++ {
			if err := c.ch.Publish("", c.q.Name, false, false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg,
				}); err != nil {
				continue
			}
			c.log.Debugf("sent notification on %d: %s", event.ID, event.Title)
			break
		}
		c.log.Warn("failed to publish a notification: ")
	}
}

func (c *Client) ConsumeAndSend(ctx context.Context, sender func(context.Context, []byte)) error {
	c.busy <- struct{}{}
	messages, err := c.ch.Consume(c.q.Name, "sender", true, false, false, false, nil)
	if err != nil {
		return err
	}
	<-c.busy
	for msg := range messages {
		sender(ctx, msg.Body)
	}
	c.busy <- struct{}{}
	return nil
}
