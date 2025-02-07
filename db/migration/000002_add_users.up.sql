CREATE TABLE "users" (
  -- "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "full_name" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");