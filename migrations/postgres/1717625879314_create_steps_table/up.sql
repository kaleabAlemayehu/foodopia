CREATE TABLE steps (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    step_order INTEGER NOT NULL
);
