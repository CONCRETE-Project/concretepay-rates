package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eabz/concretepay-rates/models"
	"github.com/shopspring/decimal"
)

var coinListCacheTimefram int64 = 24 * 60 * 60 // Every day

type coinList struct {
	data map[string]string
	lock sync.RWMutex
}

type CoinGeckoService struct {
	baseURL            string
	coinListCache      coinList
	lastCoinListUpdate int64
}

func NewCoinGeckoService() *CoinGeckoService {
	s := &CoinGeckoService{
		baseURL: "https://api.coingecko.com/api/v3/",
		coinListCache: coinList{
			data: make(map[string]string),
		},
		lastCoinListUpdate: 0,
	}
	err := s.fetchCoinList()
	if err != nil {
		panic(err)
	}
	return s
}

func (s *CoinGeckoService) GetCoinRates(tag string) (models.RatesResponse, error) {
	id, err := s.getCoinID(tag)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	res, err := client.Get(s.baseURL + "coins/" + id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var rates models.CoinGeckoCoinInfo
	err = json.Unmarshal(contents, &rates)
	if err != nil {
		return nil, err
	}
	var ratesResponse models.RatesResponse
	for code, value := range rates.MarketData.CurrentPrice {
		newRate := models.Rate{
			Code: strings.ToUpper(code),
			Rate: decimal.NewFromFloat(value),
		}
		ratesResponse = append(ratesResponse, newRate)
	}
	return ratesResponse, nil
}

func (s *CoinGeckoService) getCoinID(tag string) (string, error) {
	defer s.coinListCache.lock.Unlock()
	currTime := time.Now().Unix()
	if currTime > s.lastCoinListUpdate+coinListCacheTimefram {
		err := s.fetchCoinList()
		if err != nil {
			return "", err
		}
		s.lastCoinListUpdate = time.Now().Unix()
	}
	s.coinListCache.lock.Lock()
	id, ok := s.coinListCache.data[tag]
	if !ok {
		return "", errors.New("coin doesn't exists on coin gecko")
	}
	return id, nil
}

func (s *CoinGeckoService) fetchCoinList() error {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	res, err := client.Get(s.baseURL + "coins/list")
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var coinGeckoList models.CoinGeckoCoinListArray
	err = json.Unmarshal(contents, &coinGeckoList)
	if err != nil {
		return err
	}
	s.coinListCache.lock.Lock()
	for _, coin := range coinGeckoList {
		s.coinListCache.data[coin.Symbol] = coin.ID
	}
	s.coinListCache.lock.Unlock()
	return nil
}
