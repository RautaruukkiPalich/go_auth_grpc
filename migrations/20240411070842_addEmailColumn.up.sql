ALTER TABLE users
ADD email  VARCHAR NOT NULL UNIQUE;

DROP INDEX idx_slug;

CREATE INDEX IF NOT EXISTS idx_email ON users (email);


INSERT INTO apps (id, name, secret)
VALUES (1, 'test app', 'veryverysecretkey')
ON CONFLICT DO NOTHING;
