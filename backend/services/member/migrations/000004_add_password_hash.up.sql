-- Add password hash column to members table
ALTER TABLE members ADD COLUMN password_hash TEXT;
