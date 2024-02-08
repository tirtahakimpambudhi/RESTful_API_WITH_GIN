

-- +goose Up
-- +goose StatementBegin
CREATE TABLE todolist (
                          task_id SERIAL PRIMARY KEY,
                          user_id UUID REFERENCES users(id),
                          task_name VARCHAR(255) NOT NULL,
                          description TEXT,
                          due_date DATE,
                          priority INT,
                          completed BOOLEAN DEFAULT FALSE ,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_at_trigger
    BEFORE UPDATE ON todolist
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todolist;
-- +goose StatementEnd
