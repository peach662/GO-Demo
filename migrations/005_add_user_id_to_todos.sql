ALTER TABLE todos
    ADD COLUMN user_id BIGINT UNSIGNED NULL,
    ADD KEY idx_todos_user_id (user_id);