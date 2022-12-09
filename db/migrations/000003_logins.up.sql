CREATE TABLE IF NOT EXISTS "logins"(
    "@logins" bigserial NOT NULL UNIQUE,
    "id" text NOT NULL,
    "user" bigint REFERENCES users ("@users") ON DELETE CASCADE,
    "login" text NOT NULL,
    "password" text NOT NULL,
    "metainfo" text
);
CREATE UNIQUE INDEX IF NOT EXISTS "ilogins-user-id" ON "logins" USING btree ("id", "user");
