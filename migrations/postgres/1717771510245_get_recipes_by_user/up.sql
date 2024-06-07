CREATE OR REPLACE FUNCTION get_recipes_by_user(p_user_id INT)
RETURNS TABLE (
    id INT,
    user_id INT,
    title TEXT,
    description TEXT,
    category_id INT,
    prep_time INT,
    avg_rating FLOAT,
    total_likes INT,
    total_comments INT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    featured_image_url TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        r.id,
        r.user_id,
        r.title,
        r.description,
        r.category_id,
        r.prep_time,
        r.avg_rating,
        r.total_likes,
        r.total_comments,
        r.created_at,
        r.updated_at,
        r.featured_image_url
    FROM
        recipes r
    WHERE
        r.user_id = p_user_id;
END;
$$ LANGUAGE plpgsql VOLATILE;
