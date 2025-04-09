CREATE TABLE city_buildings
(
    id             uuid      not null primary key default gen_random_uuid(),
    city_id        uuid      not null,
    building       int       not null,
    level          int       not null,
    constructed_at timestamp not null             default current_timestamp
);

ALTER TABLE "city_buildings"
    ADD FOREIGN KEY ("city_id") REFERENCES "cities" ("id");

insert into city_buildings (city_id, building, level)
values ('715929a9-416a-4a7f-8711-a7cbb781e2dd', 0, 1);