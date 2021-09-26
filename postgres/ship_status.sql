-- Update timestamp trigger.
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Ship status table.
CREATE TABLE IF NOT EXISTS ship_status (
    capacity INT NOT NULL DEFAULT 100,
    allocated INT NOT NULL,
    loaded INT NOT NULL,
    available INT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Call the timestamp trigger on update.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON ship_status
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Seed data.
INSERT INTO ship_status (allocated, loaded, available)
VALUES (0, 0, 100);