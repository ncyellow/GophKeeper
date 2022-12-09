CREATE TABLE IF NOT EXISTS "users"(
                                      "@users" bigserial NOT NULL UNIQUE,
                                      "login" text NOT NULL,
                                      "password" text NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "iusers-login" ON "users" USING btree ("login")