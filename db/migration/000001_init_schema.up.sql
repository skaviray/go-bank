CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account" bigint NOT NULL,
  "to_account" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account");

CREATE INDEX ON "transfers" ("to_account");

CREATE INDEX ON "transfers" ("from_account", "to_account");

COMMENT ON COLUMN "entries"."amount" IS 'It can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'It must be positive';
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account") REFERENCES "accounts" ("id");
