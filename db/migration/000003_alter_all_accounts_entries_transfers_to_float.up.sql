ALTER TABLE accounts
ALTER COLUMN "balance" TYPE double precision;

ALTER TABLE entries
ALTER COLUMN "amount" type double precision;

alter table transfers
alter column "amount" type double precision;
