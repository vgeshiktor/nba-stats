package models

type PlayerStatsInput struct {
	PlayerID      int     `json:"player_id"`
	GameID        int     `json:"game_id"`
	Points        int     `json:"points"`
	Rebounds      int     `json:"rebounds"`
	Assists       int     `json:"assists"`
	Steals        int     `json:"steals"`
	Blocks        int     `json:"blocks"`
	Fouls         int     `json:"fouls"`
	Turnovers     int     `json:"turnovers"`
	MinutesPlayed float64 `json:"minutes_played"`
}
