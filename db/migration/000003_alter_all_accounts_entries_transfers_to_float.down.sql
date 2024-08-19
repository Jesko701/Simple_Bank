-- Down migration
ALTER TABLE accounts
ALTER COLUMN "balance" TYPE decimal(13,2);

ALTER TABLE entries
ALTER COLUMN "amount" TYPE decimal(13,2);

ALTER TABLE transfers
ALTER COLUMN "amount" TYPE decimal(13,2);