SET search_path TO jsonstatus, public;

CREATE TABLE IF NOT EXISTS data(
  key text PRIMARY KEY CONSTRAINT key_present CHECK ((key != '') IS TRUE),
  value text CONSTRAINT value_present CHECK ((value != '') IS TRUE),

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS key_idx ON data(key);
