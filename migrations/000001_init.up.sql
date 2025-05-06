CREATE TABLE IF NOT EXISTS clients
(
    id          VARCHAR PRIMARY KEY,
    api_key     VARCHAR   NOT NULL,
    tokens      INTEGER   NOT NULL,
    last_refill TIMESTAMP NOT NULL,
    capacity    INTEGER,
    per_second  INTEGER
);

CREATE INDEX IF NOT EXISTS idx_clients_api_key ON clients (api_key);