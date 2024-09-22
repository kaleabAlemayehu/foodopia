SET check_function_bodies = false;
INSERT INTO public.categories (id, name) VALUES (1, 'breakfast');
INSERT INTO public.categories (id, name) VALUES (2, 'lunch');
INSERT INTO public.categories (id, name) VALUES (3, 'dinner');
INSERT INTO public.categories (id, name) VALUES (4, 'dessert');
INSERT INTO public.categories (id, name) VALUES (5, 'snack');
INSERT INTO public.categories (id, name) VALUES (6, 'drink');
SELECT pg_catalog.setval('public.categories_id_seq', 1, false);
