CREATE TABLE "city_resources"
(
    "city_id"            uuid PRIMARY KEY NOT NULL,
    "wood"               float            not null default 0,
    "water"              float            not null default 0,
    "iron"               float            not null default 0,
    "oil"                float            not null default 0,
    "cotton"             float            not null default 0,
    "warehouse_capacity" float            not null default 2500,
    "updated_at"         timestamp        not null default current_timestamp,
    "population"         float            not null default 60,
    "max_population"     float            not null default 100
);

ALTER TABLE "city_resources"
    ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

INSERT INTO city_resources (city_id)
values ('715929a9-416a-4a7f-8711-a7cbb781e2dd');