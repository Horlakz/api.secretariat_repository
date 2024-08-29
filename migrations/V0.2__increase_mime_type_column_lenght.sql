-- Alter the files table to increase the length of the mime_type column
ALTER TABLE files
ALTER COLUMN mime_type TYPE VARCHAR(100);