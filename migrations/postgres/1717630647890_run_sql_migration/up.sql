CREATE OR REPLACE FUNCTION update_avg_rating_and_total_comments() RETURNS VOID AS $$
BEGIN
    UPDATE recipes
    SET avg_rating = (
        SELECT AVG(rating)
        FROM comments
        WHERE recipe_id = NEW.recipe_id
    ),
    total_comments = (
        SELECT COUNT(*)
        FROM comments
        WHERE recipe_id = NEW.recipe_id
    )
    WHERE id = NEW.recipe_id;
END;
$$ LANGUAGE plpgsql;
