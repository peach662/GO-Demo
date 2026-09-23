ALTER TABLE todos
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
    AFTER title;

UPDATE todos
SET status = CASE
    WHEN done = 1 THEN 'COMPLETED'
    ELSE 'PENDING'
END;

ALTER TABLE todos
    ADD CONSTRAINT chk_todos_status
    CHECK (status IN (
        'PENDING',
        'PROCESSING',
        'COMPLETED',
        'FAILED'
    ));
    