CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "sender" bigserial,
  "receiver" bigserial,
  "amount" decimal(10, 2),
  "created_at" timestamp
);

CREATE TABLE "wallets" (
  "id" bigserial PRIMARY KEY,
  "balance" numeric
);


ALTER TABLE "transactions" ADD FOREIGN KEY ("sender") REFERENCES "wallets" ("id");
