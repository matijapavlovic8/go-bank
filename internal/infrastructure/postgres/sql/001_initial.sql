-- +goose Up
-- SQL in section 'Up' is executed when the migration is applied
BEGIN;

-- Create the "User" table
CREATE TABLE IF NOT EXISTS "User" (
    ID serial PRIMARY KEY,
    FirstName text NOT NULL,
    LastName text NOT NULL,
    MemberSince timestamp NOT NULL,
    EncryptedPassword text NOT NULL,
    Role text NOT NULL
);
-- +migrate StatementBegin
-- Create the random_int function
CREATE OR REPLACE FUNCTION random_int(min_val int, max_val int)
    RETURNS int AS 'BEGIN RETURN floor(random() * (max_val - min_val + 1) + min_val);END;'
LANGUAGE plpgsql;
-- +migrate StatementEnd
-- Create the "Account" table
CREATE TABLE IF NOT EXISTS Account (
    AccountNumber serial PRIMARY KEY,
    Balance real NOT NULL,
    OwnerID int NOT NULL,
    Created timestamp NOT NULL
);
-- +migrate StatementEnd
-- Create the set_random_account_number function
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION set_random_account_number()
    RETURNS TRIGGER AS $$BEGIN NEW.AccountNumber := random_int(1000, 9999);RETURN NEW;END;$$
LANGUAGE plpgsql;

-- +migrate StatementEnd
-- Create the set_random_user_id function
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION set_random_user_id()
    RETURNS TRIGGER AS $$BEGIN NEW.ID := random_int(1000, 9999);RETURN NEW;END;$$
LANGUAGE plpgsql;
-- +migrate StatementEnd
-- Create the set_account_number trigger for the Account table
-- +migrate StatementBegin
DO $$BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_account_number' AND tgrelid = 'Account'::regclass) THEN CREATE TRIGGER set_account_number BEFORE INSERT ON Account FOR EACH ROW EXECUTE FUNCTION set_random_account_number();END IF;END $$;
-- +migrate StatementEnd
-- Create the set_user_id trigger for the "User" table
-- +migrate StatementBegin
DO $$BEGIN IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_user_id' AND tgrelid = '"User"'::regclass) THEN CREATE TRIGGER set_user_id BEFORE INSERT ON "User" FOR EACH ROW EXECUTE FUNCTION set_random_user_id();END IF;END $$;
-- +migrate StatementEnd
-- SQL in section 'Up' is executed when the migration is applied


-- +goose Down
-- SQL section 'Down' is executed when the migration is rolled back
-- (You can define the down migration if needed)
-- BEGIN;
-- ... (down migration SQL)
-- COMMIT;
