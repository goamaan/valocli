package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/goamaan/valocli/internal/core"
)

const (
	StorefrontUrl    = "https://pd.%s.a.pvp.net/store/v2/storefront/%s"
	BundleIdUrl      = "https://valorant-api.com/v1/bundles/%s"
	AgentsId         = "01bb38e1-da47-4e6a-9b3d-945fe4655707"
	ContractsId      = "f85cb6f7-33e5-4dc8-b609-ec7212301948"
	SpraysId         = "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475"
	GunBuddiesId     = "dd3bf334-87f3-40bd-b043-682a57a8dc3a"
	CardsId          = "3f296c07-64c3-494c-923b-fe692a4fa1bd"
	SkinsId          = "e7c63390-eda7-46e0-bb7a-a6abdacd2433"
	SkinVariantsId   = "3ad1b2b2-acdb-4524-852f-954a76ddae0a"
	TitlesId         = "de7caa6b-adf7-4588-bbd1-143831e786c6"
	CurrencyId       = "85ad13f7-3d1b-5128-9eb2-7cd8ee0b5741"
	KingdomCreditsId = "85ca954a-41f2-ce94-9b45-8ca3dd39a00d"
)

var SingleItemUrlMap = map[string]string{
	SkinsId:        "https://valorant-api.com/v1/weapons/skinlevels/%s",
	SkinVariantsId: "https://valorant-api.com/v1/weapons/skinchromas/%s",
	AgentsId:       "https://valorant-api.com/v1/agents/%s",
	ContractsId:    "https://valorant-api.com/v1/contracts/%s",
	SpraysId:       "https://valorant-api.com/v1/sprays/%s",
	GunBuddiesId:   "https://valorant-api.com/v1/buddies/levels/%s",
	CardsId:        "https://valorant-api.com/v1/playercards/%s",
	TitlesId:       "https://valorant-api.com/v1/playertitles/%s",
}

type Item struct {
	Item        string
	Cost        int
	DisplayIcon string
}

type NightMarketItem struct {
	Item            string
	BaseCost        int
	DiscountCost    int
	DiscountPercent int
	DisplayIcon     string
}

type Bundle struct {
	Items       []Item
	BundlePrice int
	DisplayName string
}

type StoreCliTable struct {
	Featured    []Bundle
	DailyStore  []Item
	Accessories []Item
	NightMarket []NightMarketItem
}

type ExternalApiSkinResponse struct {
	Data struct {
		DisplayName string `json:"displayName"`
		DisplayIcon string `json:"displayIcon"`
	} `json:"data"`
}

type StorefrontResponse struct {
	FeaturedBundle struct {
		Bundle struct {
			ID          string `json:"ID"`
			DataAssetID string `json:"DataAssetID"`
			CurrencyID  string `json:"CurrencyID"`
			Items       []struct {
				Item struct {
					ItemTypeID string `json:"ItemTypeID"`
					ItemID     string `json:"ItemID"`
					Quantity   int    `json:"Quantity"`
				} `json:"Item"`
				BasePrice       int     `json:"BasePrice"`
				CurrencyID      string  `json:"CurrencyID"`
				DiscountPercent float64 `json:"DiscountPercent"`
				DiscountedPrice float64 `json:"DiscountedPrice"`
				IsPromoItem     bool    `json:"IsPromoItem"`
			} `json:"Items"`
		} `json:"Bundle"`
		Bundles []struct {
			ID          string `json:"ID"`
			DataAssetID string `json:"DataAssetID"`
			CurrencyID  string `json:"CurrencyID"`
			Items       []struct {
				Item struct {
					ItemTypeID string `json:"ItemTypeID"`
					ItemID     string `json:"ItemID"`
					Quantity   int    `json:"Quantity"`
				} `json:"Item"`
				BasePrice       int     `json:"BasePrice"`
				CurrencyID      string  `json:"CurrencyID"`
				DiscountPercent float64 `json:"DiscountPercent"`
				DiscountedPrice float64 `json:"DiscountedPrice"`
				IsPromoItem     bool    `json:"IsPromoItem"`
			} `json:"Items"`
			TotalDiscountedCost        map[string]int `json:"TotalDiscountedCost"`
			DurationRemainingInSeconds int            `json:"DurationRemainingInSeconds"`
		} `json:"Bundles"`
		BundleRemainingDurationInSeconds int `json:"BundleRemainingDurationInSeconds"`
	} `json:"FeaturedBundle"`
	SkinsPanelLayout struct {
		SingleItemOffers      []string `json:"SingleItemOffers"`
		SingleItemStoreOffers []struct {
			OfferID          string         `json:"OfferID"`
			IsDirectPurchase bool           `json:"IsDirectPurchase"`
			StartDate        string         `json:"StartDate"`
			Cost             map[string]int `json:"Cost"`
			Rewards          []struct {
				ItemTypeID string `json:"ItemTypeID"`
				ItemID     string `json:"ItemID"`
				Quantity   int    `json:"Quantity"`
			} `json:"Rewards"`
		} `json:"SingleItemStoreOffers"`
		SingleItemOffersRemainingDurationInSeconds int `json:"SingleItemOffersRemainingDurationInSeconds"`
	} `json:"SkinsPanelLayout"`
	UpgradeCurrencyStore struct {
		UpgradeCurrencyOffers []struct {
			OfferID          string `json:"OfferID"`
			StorefrontItemID string `json:"StorefrontItemID"`
			Offer            struct {
				OfferID          string         `json:"OfferID"`
				IsDirectPurchase bool           `json:"IsDirectPurchase"`
				StartDate        string         `json:"StartDate"`
				Cost             map[string]int `json:"Cost"`
				Rewards          []struct {
					ItemTypeID string `json:"ItemTypeID"`
					ItemID     string `json:"ItemID"`
					Quantity   int    `json:"Quantity"`
				} `json:"Rewards"`
			} `json:"Offer"`
		} `json:"UpgradeCurrencyOffers"`
	} `json:"UpgradeCurrencyStore"`
	BonusStore *struct {
		BonusStoreOffers []struct {
			BonusOfferID string `json:"BonusOfferID"`
			Offer        struct {
				OfferID          string         `json:"OfferID"`
				IsDirectPurchase bool           `json:"IsDirectPurchase"`
				StartDate        string         `json:"StartDate"`
				Cost             map[string]int `json:"Cost"`
				Rewards          []struct {
					ItemTypeID string `json:"ItemTypeID"`
					ItemID     string `json:"ItemID"`
					Quantity   int    `json:"Quantity"`
				} `json:"Rewards"`
			} `json:"Offer"`
			DiscountPercent float64        `json:"DiscountPercent"`
			DiscountCosts   map[string]int `json:"DiscountCosts"`
			IsSeen          bool           `json:"IsSeen"`
		} `json:"BonusStoreOffers"`
		BonusStoreRemainingDurationInSeconds int `json:"BonusStoreRemainingDurationInSeconds"`
	} `json:"BonusStore,omitempty"`
	AccessoryStore struct {
		AccessoryStoreOffers []struct {
			Offer struct {
				OfferID          string         `json:"OfferID"`
				IsDirectPurchase bool           `json:"IsDirectPurchase"`
				StartDate        string         `json:"StartDate"`
				Cost             map[string]int `json:"Cost"`
				Rewards          []struct {
					ItemTypeID string `json:"ItemTypeID"`
					ItemID     string `json:"ItemID"`
					Quantity   int    `json:"Quantity"`
				} `json:"Rewards"`
			} `json:"Offer"`
		} `json:"AccessoryStoreOffers"`
	} `json:"AccessoryStore"`
}

func GetStorefront(c *core.Client) error {
	url := fmt.Sprintf(StorefrontUrl, c.Region, c.AuthData.UserId)
	req, err := c.RequestWithAuth("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	storefrontBody := new(StorefrontResponse)
	if err = json.NewDecoder(res.Body).Decode(&storefrontBody); err != nil {
		return err
	}

	storeCliTable := &StoreCliTable{}

	err = FetchStores(storefrontBody, storeCliTable)

	if err != nil {
		log.Fatalf("Error in fetching store items from external API:  %s", err)
	}

	return nil
}

func FetchStores(s *StorefrontResponse, table *StoreCliTable) error {
	// Daily store
	log.Println("Fetching Daily Store...")
	for _, offer := range s.SkinsPanelLayout.SingleItemStoreOffers {
		itemId := offer.Rewards[0].ItemID
		requestUrl := SingleItemUrlMap[offer.Rewards[0].ItemTypeID]
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(requestUrl, itemId), nil)
		if err != nil {
			log.Fatalf("Error creating http request to external api: %s", err)
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		responseBody := new(ExternalApiSkinResponse)
		if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatalf("Error decoding json from external api response: %s", err)
			return err
		}

		table.DailyStore = append(table.DailyStore, Item{Cost: offer.Cost[CurrencyId], Item: responseBody.Data.DisplayName, DisplayIcon: responseBody.Data.DisplayIcon})
		res.Body.Close()
	}

	// Featured bundles
	log.Println("Fetching Featured Store...")
	for _, featuredBundle := range s.FeaturedBundle.Bundles {
		bundle := new(Bundle)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(BundleIdUrl, featuredBundle.DataAssetID), nil)
		if err != nil {
			log.Fatalf("Error creating http request for bundle: %s", err)
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		responseBody := new(ExternalApiSkinResponse)
		if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatalf("Error decoding json from external api response: %s", err)
			return err
		}

		bundle.BundlePrice = featuredBundle.TotalDiscountedCost[CurrencyId]
		bundle.DisplayName = responseBody.Data.DisplayName
		for _, item := range featuredBundle.Items {
			requestUrl := SingleItemUrlMap[item.Item.ItemTypeID]
			itemId := item.Item.ItemID
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(requestUrl, itemId), nil)
			if err != nil {
				log.Fatalf("Error creating http request to external api: %s", err)
				return err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			responseBody := new(ExternalApiSkinResponse)
			if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
				log.Fatalf("Error decoding json from external api response: %s", err)
				return err
			}

			bundle.Items = append(bundle.Items, Item{Cost: item.BasePrice, Item: responseBody.Data.DisplayName, DisplayIcon: responseBody.Data.DisplayIcon})
			res.Body.Close()
		}

		res.Body.Close()
		table.Featured = append(table.Featured, *bundle)
	}

	// Night market
	if s.BonusStore != nil {
		log.Println("Fetching Night Market...")
		for _, offer := range s.BonusStore.BonusStoreOffers {
			itemId := offer.Offer.Rewards[0].ItemID
			requestUrl := SingleItemUrlMap[offer.Offer.Rewards[0].ItemTypeID]
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(requestUrl, itemId), nil)
			if err != nil {
				log.Fatalf("Error creating http request to external api: %s", err)
				return err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			responseBody := new(ExternalApiSkinResponse)
			if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
				log.Fatalf("Error decoding json from external api response: %s", err)
				return err
			}

			table.NightMarket = append(table.NightMarket,
				NightMarketItem{BaseCost: offer.Offer.Cost[CurrencyId],
					Item:            responseBody.Data.DisplayName,
					DiscountCost:    offer.DiscountCosts[CurrencyId],
					DiscountPercent: int(offer.DiscountPercent),
					DisplayIcon:     responseBody.Data.DisplayIcon})

			res.Body.Close()
		}
	}

	// Accessory Store
	log.Println("Fetching Accessories Store")
	for _, offer := range s.AccessoryStore.AccessoryStoreOffers {
		itemId := offer.Offer.Rewards[0].ItemID
		requestUrl := SingleItemUrlMap[offer.Offer.Rewards[0].ItemTypeID]
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(requestUrl, itemId), nil)
		if err != nil {
			log.Fatalf("Error creating http request to external api: %s", err)
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		responseBody := new(ExternalApiSkinResponse)
		if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatalf("Error decoding json from external api response: %s", err)
			return err
		}

		table.Accessories = append(table.Accessories,
			Item{
				Item:        responseBody.Data.DisplayName,
				Cost:        offer.Offer.Cost[KingdomCreditsId],
				DisplayIcon: responseBody.Data.DisplayIcon})

		res.Body.Close()
	}

	PrintStore(table)

	return nil
}

func PrintStore(table *StoreCliTable) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug|tabwriter.TabIndent)
	fmt.Fprintln(w, "💰 Daily store 💰")
	fmt.Fprintln(w, "Skin\tPrice\tImage Link")
	for _, item := range table.DailyStore {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%s", item.Item, item.Cost, item.DisplayIcon))
	}
	fmt.Fprintln(w, "★★★★★★★★★★★★★★★★")
	fmt.Fprintln(w, "💰 Featured store 💰")
	for _, bundle := range table.Featured {
		fmt.Fprintln(w, bundle.DisplayName)
		fmt.Fprintln(w, fmt.Sprintf("%sVP", strconv.Itoa(bundle.BundlePrice)))
		fmt.Fprintln(w, "Skin\tPrice\tImage Link")
		for _, item := range bundle.Items {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%s", item.Item, item.Cost, item.DisplayIcon))
		}
	}
	fmt.Fprintln(w, "★★★★★★★★★★★★★★★★")
	fmt.Fprintln(w, "🌟 Night Market 🌟")
	fmt.Fprintln(w, "Skin\tBase Price\tDiscount Price\tDiscount Percent\tImage Link")
	for _, item := range table.NightMarket {
		fmt.Fprintln(w,
			fmt.Sprintf("%s\t%d\t%d\t%d\t%s",
				item.Item,
				item.BaseCost,
				item.DiscountCost,
				item.DiscountPercent,
				item.DisplayIcon))
	}
	fmt.Fprintln(w, "★★★★★★★★★★★★★★★★")
	fmt.Fprintln(w, "🌟 Accessories store 🌟")
	fmt.Fprintln(w, "Item\tPrice\tImage Link")
	for _, item := range table.Accessories {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%s", item.Item, item.Cost, item.DisplayIcon))
	}
	w.Flush()
}
