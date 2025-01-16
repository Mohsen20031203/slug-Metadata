package main

import (
	"encoding/json"
	"fmt"
	"hello/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ApiResponse struct {
	Status interface{}         `json:"status"`
	Data   map[string]CoinInfo `json:"data"`
}

type CoinInfo struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	Symbol          string      `json:"symbol"`
	Category        string      `json:"category"`
	Description     string      `json:"description"`
	Slug            string      `json:"slug"`
	Logo            string      `json:"logo"`
	Subreddit       string      `json:"subreddit"`
	Notice          string      `json:"notice"`
	Tags            []string    `json:"tags"`
	TagNames        []string    `json:"tag-names"`
	TagGroups       []string    `json:"tag-groups"`
	Urls            interface{} `json:"urls"`
	Platform        interface{} `json:"platform"`
	DateAdded       time.Time   `json:"date_added"`
	TwitterUsername string      `json:"twitter_username"`
	IsHidden        int         `json:"is_hidden"`
	DateLaunched    time.Time   `json:"date_launched"`
	ContractAddress []struct {
		ContractAddress string `json:"contract_address"`
		Platform        struct {
			Name string `json:"name"`
			Coin struct {
				ID     int    `json:"id,string"`
				Name   string `json:"name"`
				Symbol string `json:"symbol"`
				Slug   string `json:"slug"`
			} `json:"coin"`
		} `json:"platform"`
	} `json:"contract_address"`
	SelfReportedCirculatingSupply interface{} `json:"self_reported_circulating_supply"`
	SelfReportedTags              interface{} `json:"self_reported_tags"`
	SelfReportedMarketCap         interface{} `json:"self_reported_market_cap"`
	InfiniteSupply                bool        `json:"infinite_supply"`
}

type CoinDetails struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	Symbol          string            `json:"symbol"`
	Slug            string            `json:"slug"`
	ContractAddress map[string]string `json:"contractAddress"`
}

func getMapURL(baseURL, endpoint string, start int) string {
	if start == 0 {
		return baseURL + endpoint
	}
	return fmt.Sprintf("%s%s&start=%d", baseURL, endpoint, start)
}

func main() {

	ch := []Chain{
		//Blast,
		Ethereum,
		BscChain,
		Base,
		Arbitrum,
		Polygon,
		Optimism,
		Avalanche,
		//CronosChain,
	}
	file2, err := os.OpenFile("Result.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening/creating file:", err)
		return
	}
	defer file2.Close()

	dataCollection, err := GetCMCContract()
	if err != nil {
		log.Fatalf("failed to fetch contract data: %v", err)
	}

	token, err := Get1InchTokens(ch[0])
	if err != nil {
		panic(err)
	}
	_ = token

	err = utils.WriteFile(file2, dataCollection)
	if err != nil {
		panic("Write in file")
	}

}

/*
   {
     "id": 1,
     "rank": 1,
     "name": "Bitcoin",
     "symbol": "BTC",
     "slug": "bitcoin",
     "is_active": 1,
     "status": 1,
     "first_historical_data": "2010-07-13T00:05:00Z",
     "last_historical_data": "2024-12-28T10:45:00Z",
     "platform": null,
     "cotracts":
     {
       "etherium":"0x21123123"
     }
*/

/*

 */

type Chain int

const (
	Ethereum Chain = iota
	Avalanche
	Base
	//Blast
	Arbitrum
	Polygon
	Optimism
	BscChain
	//CronosChain
)

func (c Chain) GetOneInchChainId() int {
	if c == Ethereum {
		return 1
	}
	if c == Avalanche {
		return 43114
	}
	if c == Base {
		return 8453
	}

	//if c == Blast {
	//	return 81457
	//}

	if c == Arbitrum {
		return 42161
	}
	if c == Polygon {
		return 137
	}
	if c == Optimism {
		return 10
	}
	if c == BscChain {
		return 56
	}
	//if c == CronosChain {
	//	return 25
	//}

	return 0
}

func (c Chain) GetCMCName() string {
	if c == Ethereum {
		return "Ethereum"
	}
	if c == Avalanche {
		return "Avalanche"
	}
	if c == Base {
		return "Base"
	}
	//if c == Blast {
	//	return ?
	//}
	if c == Arbitrum {
		return "Arbitrum"
	}
	if c == Polygon {
		return "Polygon"
	}
	if c == Optimism {
		return "Optimism"
	}
	if c == BscChain {
		return "BSC"
	}

	//if c == CronosChain {
	//	return ?
	//}
	return ""
}

func Get1InchTokens(c Chain) (utils.TokenData, error) {

	url := fmt.Sprintf("https://api.vultisig.com/1inch/swap/v6.0/%v/tokens", c.GetOneInchChainId())
	resp, err := http.Get(url)
	if err != nil {
		return utils.TokenData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.TokenData{}, err
	}

	var tokenData utils.TokenData
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return utils.TokenData{}, err
	}

	return tokenData, nil
}

const (
	baseURL      = "https://api.vultisig.com/cmc/v1/cryptocurrency"
	mapEndpoint  = "/map?sort=id&limit=5000"
	infoEndpoint = "/info?id="
)

func GetCMCContract() ([]CoinDetails, error) {

	// fetch getMap

	//

	coinDataResponse, err := getAllCoinData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coins info from API at URL '%s': %w", (baseURL + mapEndpoint), err)

	}

	// ProcessCoinDetails
	// GenerateCoinDetails
	coinIds := make([]int, 0)
	for _, c := range coinDataResponse.Data {
		coinIds = append(coinIds, c.ID)
	}
	return ProcessCoinDetails(coinIds)

}

func ProcessCoinDetails(coinIds []int) ([]CoinDetails, error) {
	var dataCollection []CoinDetails

	fullInfoURL := baseURL + infoEndpoint
	idBatch := make([]string, 0)

	for i := 0; i < len(coinIds); i++ {

		idBatch = append(idBatch, fmt.Sprintf("%d", coinIds[i]))
		if len(idBatch) == 1000 || i == len(coinIds)-1 {

			var coinsInfoResponse ApiResponse
			fullInfoURL = fullInfoURL + strings.Join(idBatch, ",")

			coinsInfoResponse, err := utils.GetApi(coinsInfoResponse, fullInfoURL)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch coins info from API at URL '%s': %w", fullInfoURL, err)
			}

			fullInfoURL = baseURL + infoEndpoint
			idBatch = nil

			for _, ci := range coinsInfoResponse.Data {

				contractAddressMap := make(map[string]string)
				for _, c := range ci.ContractAddress {
					contractAddressMap[c.Platform.Coin.Slug] = c.ContractAddress
				}

				coinDetails := CoinDetails{
					Name:            ci.Name,
					ID:              ci.ID,
					Symbol:          ci.Symbol,
					Slug:            ci.Slug,
					ContractAddress: contractAddressMap,
				}

				dataCollection = append(dataCollection, coinDetails)

			}

		}

	}

	return dataCollection, nil
}

func getAllCoinData() (utils.AutoGenerated, error) {
	var apiResponse utils.AutoGenerated

	for i := 0; ; i += 5000 {
		fullInfoURL := getMapURL(baseURL, mapEndpoint, i)

		tempResponse, err := utils.GetApi(apiResponse, fullInfoURL)
		if err != nil {
			return utils.AutoGenerated{}, fmt.Errorf("failed to fetch coins info from API at URL '%s': %w", fullInfoURL, err)
		}

		if len(tempResponse.Data) == 0 {
			break
		}

		if apiResponse.Data == nil {
			apiResponse.Data = tempResponse.Data
		} else {
			for _, value := range tempResponse.Data {

				apiResponse.Data = append(apiResponse.Data, value)

			}
		}
	}

	return apiResponse, nil
}

// Get1InchTokens(c chain) ([] tokens, err)
// GetCMCContract() ([] , err)
// main
