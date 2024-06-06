CREATE OR REPLACE FUNCTION update_avg_rating_and_total_comments() RETURNS TRIGGER AS $$
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
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_avg_rating_and_total_comments
AFTER INSERT OR UPDATE OR DELETE ON comments
FOR EACH ROW
EXECUTE PROCEDURE update_avg_rating_and_total_comments();
