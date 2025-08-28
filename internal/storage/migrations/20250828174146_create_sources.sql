-- +goose Up
-- +goose StatementBegin
CREATE TABLE sources(
    id SERIAL PRIMARY KEY,
    
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
