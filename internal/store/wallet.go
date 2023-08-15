package store

import (
	"encoding/json"
	"fmt"

	"github.com/goamaan/valocli/internal/core"
)

const (
	WalletUrl         = "https://pd.%s.a.pvp.net/store/v1/wallet/%s"
	ValorantPointsId  = "85ad13f7-3d1b-5128-9eb2-7cd8ee0b5741"
	KingdomCreditsId  = "85ca954a-41f2-ce94-9b45-8ca3dd39a00d"
	RadianitePointsId = "59aa87c-4cbf-517a-5983-6e81511be9b7"
	FreeAgentsId      = "f08d4ae3-939c-4576-ab26-09ce1f23bb37"
)

type WalletResponse struct {
	Balances map[string]int `json:"Balances"`
}

func GetWallet(c *core.Client) error {
	url := fmt.Sprintf(WalletUrl, c.Region, c.AuthData.UserId)
	req, err := c.RequestWithAuth("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	walletBody := new(WalletResponse)
	if err = json.NewDecoder(res.Body).Decode(&walletBody); err != nil {
		return err
	}
	PrintWallet(walletBody)
	return nil
}

func PrintWallet(wallet *WalletResponse) {
	fmt.Println("💵 Balances 💵")
	fmt.Printf("Valorant Points (VP) - %d\n", wallet.Balances[ValorantPointsId])
	fmt.Printf("Radianite Points (RP) - %d\n", wallet.Balances[RadianitePointsId])
	fmt.Printf("Kingdomt Credits - %d\n", wallet.Balances[KingdomCreditsId])
	fmt.Printf("Free Agents - %d\n", wallet.Balances[FreeAgentsId])
}
