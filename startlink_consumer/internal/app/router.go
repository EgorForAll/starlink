package app

import (
	"context"
	"errors"
	"strconv"

	"starlink_consumer/domain/users"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

const retryCountHeader = "retry-count"

// Alerter вызывается когда сообщение попадает в DLQ
type Alerter interface {
	Alert(ctx context.Context, msg kafka.Message, err error)
}

type Router struct {
	writer      *kafka.Writer
	alerter     Alerter
	retryTopics []string
	dlqTopic    string
	logger      zerolog.Logger
}

func NewRouter(writer *kafka.Writer, alerter Alerter, retryTopics []string, dlqTopic string, logger zerolog.Logger) *Router {
	return &Router{
		writer:      writer,
		alerter:     alerter,
		retryTopics: retryTopics,
		dlqTopic:    dlqTopic,
		logger:      logger,
	}
}

// Route направляет сообщение в retry-топик или DLQ в зависимости от типа ошибки и числа попыток
func (r *Router) Route(ctx context.Context, msg kafka.Message, err error) {
	retryCount := parseRetryCount(msg.Headers)

	var nonRetryable *users.NonRetryableError
	if !errors.As(err, &nonRetryable) && retryCount < len(r.retryTopics) {
		nextTopic := r.retryTopics[retryCount]
		r.logger.Warn().Err(err).Str("next_topic", nextTopic).Int("retry", retryCount+1).Msg("router: sending to retry")
		r.send(ctx, msg, nextTopic, retryCount+1)
		return
	}

	r.logger.Error().Err(err).Str("topic", r.dlqTopic).Msg("router: sending to DLQ")
	r.send(ctx, msg, r.dlqTopic, retryCount)
	r.alerter.Alert(ctx, msg, err)
}

func (r *Router) send(ctx context.Context, msg kafka.Message, topic string, retryCount int) {
	out := kafka.Message{
		Topic:   topic,
		Value:   msg.Value,
		Headers: setRetryCount(msg.Headers, retryCount),
	}
	if err := r.writer.WriteMessages(ctx, out); err != nil {
		r.logger.Error().Err(err).Str("topic", topic).Msg("router: failed to write message")
	}
}

func parseRetryCount(headers []kafka.Header) int {
	for _, h := range headers {
		if h.Key == retryCountHeader {
			n, _ := strconv.Atoi(string(h.Value))
			return n
		}
	}
	return 0
}

func setRetryCount(headers []kafka.Header, count int) []kafka.Header {
	out := make([]kafka.Header, 0, len(headers)+1)
	for _, h := range headers {
		if h.Key != retryCountHeader {
			out = append(out, h)
		}
	}
	return append(out, kafka.Header{
		Key:   retryCountHeader,
		Value: []byte(strconv.Itoa(count)),
	})
}
