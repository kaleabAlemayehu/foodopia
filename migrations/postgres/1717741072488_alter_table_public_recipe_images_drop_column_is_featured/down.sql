alter table "public"."recipe_images" alter column "is_featured" set default false;
alter table "public"."recipe_images" alter column "is_featured" drop not null;
alter table "public"."recipe_images" add column "is_featured" bool;
