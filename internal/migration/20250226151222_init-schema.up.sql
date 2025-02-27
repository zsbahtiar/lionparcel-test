-- Migration Up: init-schema
BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), 
  username VARCHAR(100) NOT NULL UNIQUE, 
  email VARCHAR(255) NOT NULL UNIQUE, 
  password_hash VARCHAR(255) NOT NULL, 
  -- currently set by bool, for fast development
  is_admin BOOLEAN NOT NULL DEFAULT FALSE, 
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS movies (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), 
  title VARCHAR(255) NOT NULL, 
  description TEXT, 
  -- in seconds
  duration INTEGER NOT NULL, 
  link TEXT NOT NULL, 
  -- store genres as JSON array, for fast development
  genres jsonb NOT NULL, 
  -- store artists as JSON array, for fast development
  artists jsonb NOT NULL, 
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_movie_views (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), 
  movie_id uuid NOT NULL REFERENCES movies(id),
  -- Can be NULL for anonymous users
  user_id uuid REFERENCES users(id),
  duration_watched numeric(5, 2) NOT NULL, 
  viewed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, movie_id)
);
CREATE INDEX idx_views_movie_id ON user_movie_views(movie_id);
CREATE TABLE IF NOT EXISTS votes (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), 
  user_id uuid NOT NULL REFERENCES users(id), 
  movie_id uuid NOT NULL REFERENCES movies(id), 
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
  UNIQUE(user_id, movie_id)
);
CREATE INDEX idx_votes_movie_id ON votes(movie_id);
COMMIT;
