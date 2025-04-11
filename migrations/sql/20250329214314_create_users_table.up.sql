CREATE TABLE "users"
(
    "id"             uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    "external_id"    bigint,
    "status"         string                    DEFAULT 'pending',
    "username"       string           NOT NULL,
    "password"   string NULL,
    "email"      string NOT NULL UNIQUE,
    "last_login"     timestamp,
    "last_login_ip"  inet,
    "language"       string                    DEFAULT 'en',
    "ban_expiration" timestamp,
    "ban_reason"     string,
    "created_at"     timestamp                 DEFAULT (now()),
    "updated_at" timestamp,
    "deleted_at"     timestamp
);

INSERT INTO public.users (id, username, password, email, status)
values ('c522b16b-a157-410d-8a3e-6fc64d84c17d', 'admin', '$2a$12$vgkCTdAigQmJM7MGiFV06.TJ2Lesx1tsIiC4Ehh/XbKKoeakE0wO.',
        'admin@galvanico.com', 'active');
