package models

type ValoTracker struct {
	ID string `json:"id"`
	// get account
	Puuid      string `json:"puuid"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	AccountLvl string `json:"acc_lvl"`
	CardsURL   string `json:"cards_url"`
	// get MMR
	CurrentTier string `json:"current_tier"`
	HighestRank string `json:"highest_rank"`
	// get match history
	LatestMatch []MatchRecord `json:"latest_match"`
}

type KillStats struct {
	// get match history
	ID      string `json:"id"`
	Kills   int    `json:"kills"`
	Deaths  int    `json:"deaths"`
	Assists int    `json:"assists"`
}

type MatchRecord struct {
	// get match history
	ID         string    `json:"id"`
	MatchID    string    `json:"match_id"`
	MapName    string    `json:"map_name"`
	Character  string    `json:"character"`
	TeamColor  string    `json:"team_color"`
	Finalscore string    `json:"finalscore"`
	Stats      KillStats `json:"stats"`
	//get mmr history
	MMRChange string `json:"mmr_change"`
}
