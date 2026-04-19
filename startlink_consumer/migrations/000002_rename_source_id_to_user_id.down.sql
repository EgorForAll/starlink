ALTER TABLE received_users RENAME COLUMN user_id TO source_id;
ALTER INDEX idx_received_users_user_id RENAME TO idx_received_users_source_id;
