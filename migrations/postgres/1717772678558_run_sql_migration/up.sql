-- Create or replace the function with STABLE keyword
CREATE OR REPLACE FUNCTION get_recipes_for_user(user_row users)
RETURNS SETOF recipes AS $$
BEGIN
    RETURN QUERY
    SELECT
        r.*
    FROM
        recipes r
    WHERE
        r.user_id = user_row.id;
END;
$$ LANGUAGE plpgsql STABLE;

-- Tracking the computed field (assuming done via Hasura Console);
