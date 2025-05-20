CREATE TABLE "cities"
(
    "id"       uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    "name"     string           NOT NULL,
    position_x int              not null,
    position_y int              not null,
    created_at timestamp not null default current_timestamp,
    INDEX city_position_idx (position_x, position_y)
);

INSERT INTO cities (id, name, position_x, position_y)
values ('715929a9-416a-4a7f-8711-a7cbb781e2dd', 'default', 50, 50);