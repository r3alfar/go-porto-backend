package models

// 3rd party Raw Response in Struct
type AccountDetail struct {
	Data struct {
		AccountLevel int `json:"account_level"`
		Card         struct {
			ID    string `json:"id"`
			Large string `json:"large"`
			Small string `json:"small"`
			Wide  string `json:"wide"`
		} `json:"card"`
		LastUpdate    string `json:"last_update"`
		LastUpdateRaw int    `json:"last_update_raw"`
		Name          string `json:"name"`
		Puuid         string `json:"puuid"`
		Region        string `json:"region"`
		Tag           string `json:"tag"`
	} `json:"data"`
	Status int `json:"status"`
}

type Matchlist struct {
	Status int `json:"status"`
	Data   []struct {
		Metadata struct {
			Map        string `json:"map"`
			GameLength int    `json:"game_length"`
			GameStart  int    `json:"game_start"`
			MatchId    string `json:"matchid"`
			Mode       string `json:"mode"`
			Region     string `json:"region"`
			Cluster    string `json:"cluster"`
		} `json:"metadata"`
		Players struct {
			AllPlayers []struct {
				PUUID         string `json:"puuid"`
				Name          string `json:"name"`
				Tag           string `json:"tag"`
				Team          string `json:"team"`
				Character     string `json:"character"`
				CurrentTier   string `json:"currenttier_patched"`
				CurrentTierId int    `json:"currenttier"`
				Assets        struct {
					Card struct {
						Wide string `json:"wide"`
					} `json:"card"`
					Agent struct {
						Killfeed string `json:"killfeed"`
					} `json:"agent"`
				} `json:"assets"`
				Stats struct {
					Score     int `json:"score"`
					Kills     int `json:"kills"`
					Deaths    int `json:"deaths"`
					Assists   int `json:"assists"`
					Bodyshots int `json:"bodyshots"`
					Headshots int `json:"headshots"`
					Legshots  int `json:"legshots"`
				}
			} `json:"all_players"`
		} `json:"players"`
		Teams struct {
			Red  TeamStruct `json:"red"`
			Blue TeamStruct `json:"blue"`
		} `json:"teams"`
	} `json:"data"`
}

type TeamStruct struct {
	HasWon     bool `json:"has_won"`
	RoundsWon  int  `json:"rounds_won"`
	RoundsLost int  `json:"rounds_lost"`
}

type MMR struct {
	Status int `json:"status"`
	Data   struct {
		Name        string `json:"name"`
		Tag         string `json:"Tag"`
		HighestRank struct {
			Tier        int    `json:"tier"`
			PatchedTier string `json:"patched_tier"`
			Season      string `json:"season"`
		} `json:"highest_rank"`
		CurrentData struct {
			CurrentTierID int    `json:"currenttier"`
			CurrentTier   string `json:"currenttierpatched"`
			Images        struct {
				TierIconURL string `json:"small"`
			} `json:"images"`
		} `json:"current_data"`
	} `json:"data"`
}

type MMRHistory struct {
	Data []struct {
		Tier      string `json:"currenttierpatched"`
		TierId    int    `json:"currenttier"`
		MatchId   string `json:"match_id"`
		MMRChange int    `json:"mmr_change_to_last_game"`
		DateRaw   int    `json:"date_raw"`
		Map       struct {
			Name  string `json:"name"`
			MapId string `json:"id"`
		} `json:"map"`
		Images struct {
			Small string `json:"small"`
		} `json:"images"`
	} `json:"data"`
}

type MapDetail struct {
	DisplayName  string `json:"displayName"`
	ListViewIcon string `json:"listViewIcon"`
}

// From DB
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
