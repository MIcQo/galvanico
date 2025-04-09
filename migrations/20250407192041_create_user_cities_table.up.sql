CREATE TABLE "user_cities"
(
    "user_id" uuid NOT NULL,
    "city_id" uuid NOT NULL,
    PRIMARY KEY ("user_id", "city_id")
);

ALTER TABLE "user_cities"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_cities"
    ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

INSERT INTO user_cities (user_id, city_id)
values ('c522b16b-a157-410d-8a3e-6fc64d84c17d', '715929a9-416a-4a7f-8711-a7cbb781e2dd');