CREATE TABLE IF NOT EXISTS "cards"(
    "@cards" bigserial NOT NULL UNIQUE,
    "id" text NOT NULL,
    "user" bigint REFERENCES users ("@users") ON DELETE CASCADE,
    "fio" text NOT NULL,
    "number" text NOT NULL,
    "date" text NOT NULL,
    "cvv" text NOT NULL,
    "metainfo" text
);
CREATE UNIQUE INDEX IF NOT EXISTS "icards-user-id" ON "cards" USING btree ("id", "user");
