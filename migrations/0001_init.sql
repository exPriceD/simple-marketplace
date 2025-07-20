CREATE TABLE IF NOT EXISTS schema_migrations
(
    id         SERIAL PRIMARY KEY,
    name       TEXT        NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users
(
    id            TEXT PRIMARY KEY,
    login         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_login_format CHECK (login ~ '^[A-Za-z0-9_]{3,32}$'),
    CONSTRAINT users_password_len CHECK (char_length(password_hash) BETWEEN 10 AND 512)
);

CREATE TABLE IF NOT EXISTS listings
(
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL,
    image_url   TEXT        NOT NULL,
    price       BIGINT      NOT NULL,
    author_id   TEXT        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT listings_title_len CHECK (char_length(title) BETWEEN 1 AND 120),
    CONSTRAINT listings_desc_len CHECK (char_length(description) BETWEEN 1 AND 2000),
    CONSTRAINT listings_image_len CHECK (char_length(image_url) BETWEEN 1 AND 1024),
    CONSTRAINT listings_price_range CHECK (price BETWEEN 0 AND 1000000000)
);
