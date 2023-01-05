-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     name TEXT,
--     email TEXT,
--     created_at timestamptz NOT NULL DEFAULT(now())
-- );

CREATE TABLE short_links (
    id SERIAL PRIMARY KEY,
    -- user_id INT NOT NULL,
    short_url TEXT NOT NULL,
    original_url TEXT NOT NULL,
    -- CONSTRAINT fk_users
    --     FOREIGN KEY(user_id)
    --         REFERENCES users(id)
);
