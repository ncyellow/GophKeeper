CREATE TABLE IF NOT EXISTS "bin_data"(
    "@bin" bigserial NOT NULL UNIQUE,
    "id" text NOT NULL,
    "user" bigint REFERENCES users ("@users") ON DELETE CASCADE,
    "content" text NOT NULL,
    "metainfo" text
);
CREATE UNIQUE INDEX IF NOT EXISTS "ibin_data-user-id" ON "bin_data" USING btree ("id", "user");