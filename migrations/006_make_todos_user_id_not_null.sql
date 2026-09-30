-- 历史 NULL 行回填到用户 1，再改为 NOT NULL。

UPDATE todos SET user_id = 1 WHERE user_id IS NULL;

ALTER TABLE todos MODIFY COLUMN user_id BIGINT UNSIGNED NOT NULL;