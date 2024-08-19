-- Drop the unique constraint of the "accounts" table
ALTER TABLE IF EXISTS accounts
DROP CONSTRAINT IF EXISTS "unique_owner_currency";

-- Drop the foreign key constraint on the "accounts" table
ALTER TABLE IF EXISTS accounts
DROP CONSTRAINT IF EXISTS "fk_owner_username";

-- Delete all rows from the "users"table
DELETE FROM users;

-- Drop the "users" table
DROP TABLE IF EXISTS users;