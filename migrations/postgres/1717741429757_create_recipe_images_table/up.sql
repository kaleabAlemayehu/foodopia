CREATE TABLE recipe_images (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL, -- URL or path to the image
    is_featured BOOLEAN DEFAULT FALSE -- To mark one image as the featured image
);
