CREATE OR REPLACE FUNCTION update_total_likes() RETURNS TRIGGER AS $$
BEGIN
    UPDATE recipes
    SET total_likes = (
        SELECT COUNT(*)
        FROM likes
        WHERE recipe_id = NEW.recipe_id
    )
    WHERE id = NEW.recipe_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_total_likes
AFTER INSERT OR DELETE ON likes
FOR EACH ROW
EXECUTE PROCEDURE update_total_likes();