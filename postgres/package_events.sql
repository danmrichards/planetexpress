CREATE TABLE IF NOT EXISTS package_events (
    event_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    package_id VARCHAR(50) NOT NULL,
    package_size INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id)
);