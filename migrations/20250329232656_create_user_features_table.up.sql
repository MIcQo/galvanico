CREATE TABLE "user_features"
(
    "user_id" uuid   NOT NULL,
    "feature" string NOT NULL,
    PRIMARY KEY ("user_id", "feature")
);

ALTER TABLE "user_features"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

INSERT INTO user_features (user_id, feature)
values ('c522b16b-a157-410d-8a3e-6fc64d84c17d', 'premium');
