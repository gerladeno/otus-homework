-- +goose Up
-- +goose StatementBegin
CREATE table events
(
    id          integer primary key,
    title       text,
    start_time  timestamp,
    duration    integer,
    description text,
    owner       integer,
    notify_time integer,
    created     timestamp default now(),
    updated     timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table events;
-- +goose StatementEnd