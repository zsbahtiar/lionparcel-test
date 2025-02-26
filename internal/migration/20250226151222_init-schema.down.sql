-- Migration Down: rollback init-schema
BEGIN;

-- Drop indexes first
DROP INDEX IF EXISTS idx_votes_movie_id;
DROP INDEX IF EXISTS idx_views_movie_id;

-- Drop tables with foreign key dependencies first
DROP TABLE IF EXISTS user_watch_movie_durations;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS user_movie_views;
DROP TABLE IF EXISTS movies;
DROP TABLE IF EXISTS users;

-- Drop the extension last
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;