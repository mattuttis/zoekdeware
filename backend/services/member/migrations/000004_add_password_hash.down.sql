-- Remove password hash column from members table
ALTER TABLE members DROP COLUMN password_hash;
