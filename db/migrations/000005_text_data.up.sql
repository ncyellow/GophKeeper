CREATE TABLE IF NOT EXISTS "text_data"(
    "@text" bigserial NOT NULL UNIQUE,
    "id" text NOT NULL,
    "user" bigint REFERENCES users ("@users") ON DELETE CASCADE,
    "content" text NOT NULL,
    "metainfo" text
);
CREATE UNIQUE INDEX IF NOT EXISTS "itext_data-user-id" ON "text_data" USING btree ("id", "user");