--
-- Seed data for the application
--

INSERT INTO "public"."users"("email","first_name","last_name","password","user_active", "is_admin", "created_at","updated_at")
VALUES
    (E'admin@example.com',E'Admin',E'User',E'$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',1,1,E'2022-03-14 00:00:00',E'2022-03-14 00:00:00');

SELECT pg_catalog.setval('public.plans_id_seq', 1, false);


SELECT pg_catalog.setval('public.user_id_seq', 2, true);


SELECT pg_catalog.setval('public.user_plans_id_seq', 1, false);

INSERT INTO "public"."plans"("plan_name","plan_amount","created_at","updated_at")
VALUES
    (E'Bronze Plan',1000,E'2022-05-12 00:00:00',E'2022-05-12 00:00:00'),
    (E'Silver Plan',2000,E'2022-05-12 00:00:00',E'2022-05-12 00:00:00'),
    (E'Gold Plan',3000,E'2022-05-12 00:00:00',E'2022-05-12 00:00:00');