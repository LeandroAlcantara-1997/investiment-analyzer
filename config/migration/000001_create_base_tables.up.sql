CREATE TYPE STATUS AS ENUM ('success', 'error');

CREATE TABLE logger(
    id UUID PRIMARY KEY,
    consultation_time TIMESTAMP,
    ip VARCHAR(50),
    query_parameters VARCHAR(100),
    status STATUS
);

