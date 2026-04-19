ALTER TABLE received_users RENAME COLUMN source_id TO user_id;
ALTER INDEX idx_received_users_source_id RENAME TO idx_received_users_user_id;
