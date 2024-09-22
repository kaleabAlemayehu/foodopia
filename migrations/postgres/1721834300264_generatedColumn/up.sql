ALTER TABLE ingredients
ADD COLUMN name_quantity VARCHAR(150) GENERATED ALWAYS AS (name || ' - ' || quantity) STORED;
