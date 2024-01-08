\connect tickets

CREATE TABLE tickets
(
    id UUID PRIMARY KEY,
    version NUMERIC NOT NULL DEFAULT 0,
    summary VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    status VARCHAR(50) NOT NULL
);

CREATE TRIGGER version_trigger
   BEFORE UPDATE ON tickets
   FOR EACH ROW EXECUTE PROCEDURE increment_version();