-- Account table
CREATE TABLE "account" (
  "account_id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

-- Entry table. This will record the top-up/withdrawal, so amount can be positive or negative
CREATE TABLE "entry" (
  "entry_id" bigserial PRIMARY KEY,
  "account_id" bigint not null ,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

-- Transfer table. This table record the transaction between 2 accounts, thus amount must be positive
CREATE TABLE "transfer" (
  "transfer_id" bigserial PRIMARY KEY,
  "from_account_id" bigint not null ,
  "to_account_id" bigint not null ,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entry" ("account_id");

CREATE INDEX ON "transfer" ("from_account_id");

CREATE INDEX ON "transfer" ("to_account_id");

CREATE INDEX ON "transfer" ("from_account_id", "to_account_id");

ALTER TABLE "entry" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("account_id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("from_account_id") REFERENCES "account" ("account_id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("to_account_id") REFERENCES "account" ("account_id");