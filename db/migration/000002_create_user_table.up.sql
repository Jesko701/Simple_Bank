CREATE TABLE IF NOT EXISTS users (
    "username" varchar PRIMARY KEY,
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL default '0001-01-01 00:00:00Z',
    "created_at" timestamptz not null default CURRENT_TIMESTAMP
);

BEGIN;
  INSERT INTO "users" (username, hashed_password, full_name, email)
  SELECT DISTINCT owner, 'temp_hash', owner, owner || '@temp.com'
  FROM "accounts"
  ON CONFLICT (username) DO NOTHING;
COMMIT;

ALTER TABLE IF EXISTS accounts
ADD CONSTRAINT "fk_owner_username" FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- Unique constraint to ensure an owner cannot have duplicate accounts with the same currency
ALTER TABLE IF EXISTS accounts
ADD CONSTRAINT "unique_owner_currency" UNIQUE ("owner", "currency");