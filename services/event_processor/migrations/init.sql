CREATE TABLE IF NOT EXISTS events
(
    id         SERIAL PRIMARY KEY,
    event_type VARCHAR(50),
    page_url   TEXT,
    timestamp  TIMESTAMP,
    user_id    VARCHAR(50)
);

select * from events;