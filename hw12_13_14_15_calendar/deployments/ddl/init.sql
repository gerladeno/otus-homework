CREATE TABLE IF NOT EXISTS events
(
    id          serial primary key,
    title       text,
    start_time  timestamp,
    duration    integer,
    description text,
    owner       integer,
    notify_time integer,
    created     timestamp default now(),
    updated     timestamp default now()
);