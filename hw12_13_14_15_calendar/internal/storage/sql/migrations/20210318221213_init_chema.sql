-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id          serial primary key,
    title       text      not null,
    start_time  timestamp not null,
    duration    integer   not null,
    description text      not null,
    owner       integer   not null,
    notify_time integer   not null,
    created     timestamp default now(),
    updated     timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table events;
-- +goose StatementEnd