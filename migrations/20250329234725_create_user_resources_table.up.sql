CREATE TABLE "user_resources"
(
    "user_id"              uuid PRIMARY KEY NOT NULL,
    "gold"                 integer          NOT NULL DEFAULT 10000,
    "train_cars"           integer          NOT NULL DEFAULT 0,
    "available_train_cars" integer          NOT NULL DEFAULT 0,
    "electricity"          float            NOT NULL DEFAULT 0,
    "max_electricity"      float            NOT NULL DEFAULT 100,
    "waste"                float            NOT NULL DEFAULT 0,
    "max_waste"            float            NOT NULL DEFAULT 100,
    "updated_at"           timestamp        NOT NULL DEFAULT now()
);

ALTER TABLE "user_resources"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE cascade ON UPDATE no action;

insert into user_resources (user_id)
values ('c522b16b-a157-410d-8a3e-6fc64d84c17d');