alter table "public"."likes" add constraint "likes_user_id_recipe_id_key" unique ("user_id", "recipe_id");
