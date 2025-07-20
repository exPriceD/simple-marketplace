CREATE INDEX IF NOT EXISTS listings_created_at_idx ON listings (created_at);
CREATE INDEX IF NOT EXISTS listings_price_idx ON listings (price);
CREATE INDEX IF NOT EXISTS listings_price_created_at_idx ON listings (price, created_at);
CREATE INDEX IF NOT EXISTS listings_author_id_idx ON listings (author_id);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_uq_idx ON users (login);