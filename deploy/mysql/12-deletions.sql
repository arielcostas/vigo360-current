-- Update the publicaciones table, and add a legal_restricted_at column with a default value of NULL
ALTER TABLE publicaciones ADD COLUMN legally_retired_at DATETIME DEFAULT NULL;