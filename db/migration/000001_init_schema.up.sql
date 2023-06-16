CREATE TABLE "account" (
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

CREATE TABLE "Transfers" (
  "id" bigserial PRIMARY KEY,
  "from_accound_id" bigint NOT NULL,
  "to_accound_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "Transfers" ("from_accound_id");

CREATE INDEX ON "Transfers" ("to_accound_id");

CREATE INDEX ON "Transfers" ("from_accound_id", "to_accound_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be postive or negative';

COMMENT ON COLUMN "Transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("from_accound_id") REFERENCES "account" ("id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("to_accound_id") REFERENCES "account" ("id");
