alter table "public"."recipe_images" alter column "image_url" drop not null;
alter table "public"."recipe_images" add column "image_url" text;
