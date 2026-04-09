CREATE TABLE IF NOT EXISTS received_users (
    id           BIGSERIAL    PRIMARY KEY,
    source_id    BIGINT       NOT NULL,          -- id пользователя из producer_db
    first_name   VARCHAR(255) NOT NULL,           -- обработанные данные (верхний регистр)
    last_name    VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL,
    received_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ  NOT NULL            -- время обработки сообщения
);

CREATE INDEX IF NOT EXISTS idx_received_users_source_id
    ON received_users (source_id);
