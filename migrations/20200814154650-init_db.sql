-- +migrate Up
CREATE TYPE "request_type" AS ENUM ('out', 'in');

CREATE TABLE "tag"
(
    "id"   SERIAL PRIMARY KEY,
    "name" varchar
);

CREATE TABLE "request"
(
    "id"         SERIAL PRIMARY KEY,
    "type"       request_type NOT NULL,
    "tg_user_id" VARCHAR      NOT NULL,
    "tag_id"     INT          NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT tag_id_fkey FOREIGN KEY ("tag_id")
        REFERENCES tag (id) MATCH SIMPLE
        ON UPDATE NO ACTION ON DELETE NO ACTION,
    UNIQUE ("type", "tg_user_id", "tag_id")
);

-- +migrate Down
DROP TABLE "tag";
DROP TABLE "request";
DROP TYPE "request_type";
