CREATE TABLE IF NOT EXISTS sc_players.players (
    "id" BIGSERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    CONSTRAINT "PK_Players" PRIMARY KEY ("id")
);