SET check_function_bodies = false;
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (81, 'Kaleab Alemayehu', 'kaleab@gmail.com', '$2a$10$o2Ps5vxVyipHCOUiYiq3/ucH2xIDdRNnH0LLipcFlpI8EPROST.Au', '2024-07-18 14:26:24.014268');
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (82, 'user', 'user@gmail.com', '$2a$10$pCQ3.hDOj.lBjvgYiz7OwOhvAqLOvBxBjwP8XdvPt9sNTqHv83ixO', '2024-07-19 12:15:40.07706');
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (83, 'test', 'test@gamil.com', '$2a$10$yeWIvbUhJ1vR96coD9paW.EloVJOnpaknat0VOF5JD/dA4xbvsEgW', '2024-07-19 13:46:06.261411');
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (94, 'testUser', 'testUser@gmail.com', '$2a$10$YEphAck1jh.rzTCQykEZXOU9mhgfd.ed3/Ce3lVaODj02a8WJt4Dm', '2024-07-21 12:31:24.364848');
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (95, 'toast', 'toast@toast.com', '$2a$10$.uLKKU1mxCHsbxpjF7OXUeil42CGdjVFCufCMzVRU.lOGMf5AL91a', '2024-07-25 08:32:27.205309');
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (98, 'kaleab', 'kaleabalemayehu04@gmail.com', '$2a$10$FF1AXmxcxml5bngy206Pp.PoAkUoGNaPs.FrAaHF7SbJxWSQuQUSi', '2024-07-25 11:30:26.603828');
INSERT INTO public.users (id, username, email, password_hash, created_at) VALUES (99, 'Michael.Tesfaye@hahu.jobs', 'Michael.Tesfaye@hahu.jobs', '$2a$10$pXLtLT2UTnmfAgRyj01w5ezIuTXvi.hehLupHsUMOK2PoB5SkteIC', '2024-09-13 21:01:04.469372');
SELECT pg_catalog.setval('public.users_id_seq', 99, true);
