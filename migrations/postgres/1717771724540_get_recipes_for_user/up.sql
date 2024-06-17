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
$$ LANGUAGE plpgsql VOLATILE;