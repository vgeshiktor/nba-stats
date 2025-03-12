-- Drop Players table
DROP TABLE IF EXISTS players CASCADE;

-- Create Players table
CREATE TABLE IF NOT EXISTS players (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    team_id TEXT NOT NULL
);

-- Drop Teams table
DROP TABLE IF EXISTS teams CASCADE;

-- Create Teams table
CREATE TABLE IF NOT EXISTS teams (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Drop Games table
DROP TABLE IF EXISTS games CASCADE;

-- Create Games table
CREATE TABLE IF NOT EXISTS games (
    id TEXT PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    home_team TEXT NOT NULL,
    away_team TEXT NOT NULL
);

-- Drop PlayerGameStats table
DROP TABLE IF EXISTS player_game_stats CASCADE;

-- Create PlayerGameStats table
CREATE TABLE IF NOT EXISTS player_game_stats (
    id TEXT PRIMARY KEY,
    player_id TEXT NOT NULL,
    game_id TEXT NOT NULL,
    points INTEGER NOT NULL,
    rebounds INTEGER NOT NULL,
    assists INTEGER NOT NULL,
    steals INTEGER NOT NULL,
    blocks INTEGER NOT NULL,
    fouls INTEGER NOT NULL,
    turnovers INTEGER NOT NULL,
    minutes_played FLOAT NOT NULL,
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (game_id) REFERENCES games(id)
);
