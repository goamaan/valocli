package core

import (
	"encoding/json"
	"net/http"
	"time"
)

type Version struct {
	Status int          `json:"status"`
	Data   ManifestData `json:"data"`
}

type ManifestData struct {
	ManifestID        string    `json:"manifestId"`
	Branch            string    `json:"branch"`
	Version           string    `json:"version"`
	BuildVersion      string    `json:"buildVersion"`
	EngineVersion     string    `json:"engineVersion"`
	RiotClientVersion string    `json:"riotClientVersion"`
	RiotClientBuild   string    `json:"riotClientBuild"`
	BuildDate         time.Time `json:"buildDate"`
}

type CompetitiveTierResponse struct {
	Status int                           `json:"status"`
	Data   []CompetitiveTierResponseData `json:"data"`
}

type CompetitiveTierResponseData struct {
	Uuid            string `json:"uuid"`
	AssetObjectName string `json:"assetObjectName"`
	Tiers           []struct {
		Tier     int    `json:"tier"`
		TierName string `json:"tierName"`
	} `json:"tiers"`
}

func GetClientVersion() (*string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://valorant-api.com/v1/version", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	versionBody := new(Version)
	if err = json.NewDecoder(res.Body).Decode(&versionBody); err != nil {
		return nil, err
	}

	return &versionBody.Data.RiotClientVersion, nil
}

func GetCompetitiveTiers() (map[int]string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://valorant-api.com/v1/competitivetiers", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	tiersBody := new(CompetitiveTierResponse)
	if err = json.NewDecoder(res.Body).Decode(&tiersBody); err != nil {
		return nil, err
	}

	tierMap := make(map[int]string)

	for _, tier := range tiersBody.Data[len(tiersBody.Data)-1].Tiers {
		tierMap[tier.Tier] = tier.TierName
	}

	return tierMap, nil
}
