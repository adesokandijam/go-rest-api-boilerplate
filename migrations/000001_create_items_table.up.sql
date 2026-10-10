CREATE TABLE IF NOT EXISTS items (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    name text NOT NULL,
    description text NOT NULL DEFAULT '',
    version integer NOT NULL DEFAULT 1
);
