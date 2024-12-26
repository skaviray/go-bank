CREATE TABLE "Accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Transfers" (
  "id" bigint PRIMARY KEY,
  "from_account" bigint NOT NULL,
  "to_account" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "Accounts" ("owner");

CREATE INDEX ON "Entries" ("account_id");

CREATE INDEX ON "Transfers" ("from_account");

CREATE INDEX ON "Transfers" ("to_account");

CREATE INDEX ON "Transfers" ("from_account", "to_account");

COMMENT ON COLUMN "Entries"."amount" IS 'It can be negative or positive';

COMMENT ON COLUMN "Transfers"."amount" IS 'It must be positive';

ALTER TABLE "Entries" ADD FOREIGN KEY ("account_id") REFERENCES "Accounts" ("id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("from_account") REFERENCES "Accounts" ("id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("to_account") REFERENCES "Accounts" ("id");
