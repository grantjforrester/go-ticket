\connect tickets

CREATE OR REPLACE FUNCTION increment_version()
  RETURNS TRIGGER
AS
$body$
BEGIN
  new.version := new.version + 1;
  RETURN new;
END;
$body$
LANGUAGE plpgsql;