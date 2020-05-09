package controller

import (
	"sort"

	"github.com/CONCRETE-Project/concretepay-rates/models"
	"github.com/CONCRETE-Project/concretepay-rates/services"
	"github.com/gin-gonic/gin"
)

type RateController struct {
	CoinGeckoService *services.CoinGeckoService
}

func (r *RateController) GetCoin(c *gin.Context) {
	tag := c.Param("coin")
	wantMap := c.Query("map")
	if tag == "" {
		c.JSON(500, gin.H{"error": "no tag specified", "status": -1})
		return
	}
	rates, err := r.CoinGeckoService.GetCoinRates(tag)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "status": -1})
		return
	}
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Code < rates[j].Code
	})
	if wantMap == "1" {
		ratesMap := make(map[string]models.Rate)
		for _, rate := range rates {
			ratesMap[rate.Code] = rate
		}
		c.JSON(200, gin.H{"status": 1, "data": ratesMap})
		return
	}
	c.JSON(200, gin.H{"status": 1, "data": rates})
	return
}
