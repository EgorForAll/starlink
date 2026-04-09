-- Таблица для паттерна Outbox.
-- Событие пишется в одной транзакции с основным INSERT (например users),
-- после чего relay-воркер читает её и отправляет в Kafka.

CREATE TABLE IF NOT EXISTS outbox (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type   VARCHAR(255) NOT NULL,
    payload      JSONB        NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ            -- NULL = ещё не обработано
);

-- Индекс только по необработанным событиям — relay использует его при поллинге.
CREATE INDEX IF NOT EXISTS idx_outbox_unprocessed
    ON outbox (created_at)
    WHERE processed_at IS NULL;
