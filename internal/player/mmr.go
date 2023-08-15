package player

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/goamaan/valocli/internal/core"
)

const (
	PlayerMMRUrl = "https://pd.%s.a.pvp.net/mmr/v1/players/%s"
)

type PlayerMMRResponse struct {
	Version                     int                     `json:"Version"`
	Subject                     string                  `json:"Subject"`
	NewPlayerExperienceFinished bool                    `json:"NewPlayerExperienceFinished"`
	QueueSkills                 map[string]QueueSkill   `json:"QueueSkills"`
	LatestCompetitiveUpdate     LatestCompetitiveUpdate `json:"LatestCompetitiveUpdate"`
	IsLeaderboardAnonymized     bool                    `json:"IsLeaderboardAnonymized"`
	IsActRankBadgeHidden        bool                    `json:"IsActRankBadgeHidden"`
}

type QueueSkill struct {
	TotalGamesNeededForRating         int                     `json:"TotalGamesNeededForRating"`
	TotalGamesNeededForLeaderboard    int                     `json:"TotalGamesNeededForLeaderboard"`
	CurrentSeasonGamesNeededForRating int                     `json:"CurrentSeasonGamesNeededForRating"`
	SeasonalInfoBySeasonID            map[string]SeasonalInfo `json:"SeasonalInfoBySeasonID"`
}

type SeasonalInfo struct {
	SeasonID                   string         `json:"SeasonID"`
	NumberOfWins               int            `json:"NumberOfWins"`
	NumberOfWinsWithPlacements int            `json:"NumberOfWinsWithPlacements"`
	NumberOfGames              int            `json:"NumberOfGames"`
	Rank                       int            `json:"Rank"`
	CapstoneWins               int            `json:"CapstoneWins"`
	LeaderboardRank            int            `json:"LeaderboardRank"`
	CompetitiveTier            int            `json:"CompetitiveTier"`
	RankedRating               int            `json:"RankedRating"`
	WinsByTier                 map[string]int `json:"WinsByTier,omitempty"`
	GamesNeededForRating       int            `json:"GamesNeededForRating"`
	TotalWinsNeededForRank     int            `json:"TotalWinsNeededForRank"`
}

type LatestCompetitiveUpdate struct {
	MatchID                      string `json:"MatchID"`
	MapID                        string `json:"MapID"`
	SeasonID                     string `json:"SeasonID"`
	MatchStartTime               int    `json:"MatchStartTime"`
	TierAfterUpdate              int    `json:"TierAfterUpdate"`
	TierBeforeUpdate             int    `json:"TierBeforeUpdate"`
	RankedRatingAfterUpdate      int    `json:"RankedRatingAfterUpdate"`
	RankedRatingBeforeUpdate     int    `json:"RankedRatingBeforeUpdate"`
	RankedRatingEarned           int    `json:"RankedRatingEarned"`
	RankedRatingPerformanceBonus int    `json:"RankedRatingPerformanceBonus"`
	CompetitiveMovement          string `json:"CompetitiveMovement"`
	AFKPenalty                   int    `json:"AFKPenalty"`
}

func GetPlayerMMR(c *core.Client) error {
	url := fmt.Sprintf(PlayerMMRUrl, c.Region, c.AuthData.UserId)
	req, err := c.RequestWithAuth("GET", url, nil)
	if err != nil {
		return err
	}
	clientVersion, err := core.GetClientVersion()

	if err != nil {
		return err
	}

	req.Header.Add("X-Riot-ClientPlatform", "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjogIldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9")
	req.Header.Add("X-Riot-ClientVersion", *clientVersion)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	playerMMRBody := new(PlayerMMRResponse)
	if err = json.NewDecoder(res.Body).Decode(&playerMMRBody); err != nil {
		return err
	}

	log.Println("player: ", playerMMRBody)

	return nil
}
