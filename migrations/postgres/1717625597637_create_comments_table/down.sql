-- Could not auto-generate a down migration.
-- Please write an appropriate down migration for the SQL below:
-- CREATE TABLE comments (
--     id SERIAL PRIMARY KEY,
--     user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
--     recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
--     comment TEXT NOT NULL,
--     rating INTEGER CHECK (rating >= 1 AND rating <= 5),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
