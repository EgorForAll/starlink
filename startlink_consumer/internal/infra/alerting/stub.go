package alerting

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

// StubAlerter имитирует отправку алерта во внешний сервис
type StubAlerter struct {
	logger zerolog.Logger
}

func NewStubAlerter(logger zerolog.Logger) *StubAlerter {
	return &StubAlerter{logger: logger}
}

func (a *StubAlerter) Alert(_ context.Context, msg kafka.Message, err error) {
	a.logger.Error().
		Err(err).
		Str("topic", msg.Topic).
		Bytes("payload", msg.Value).
		Msg("ALERT: message moved to DLQ, manual intervention required")
}
