package main

import (
	"net/http"
	"os"
	"time"

	controller "github.com/CONCRETE-Project/concretepay-rates/controllers"
	"github.com/CONCRETE-Project/concretepay-rates/services"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	App := GetApp()
	err := App.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

// GetApp is used to wrap all the additions to the GIN API.
func GetApp() *gin.Engine {
	App := gin.Default()
	App.Use(cors.Default())
	ApplyRoutes(App)
	return App
}

// ApplyRoutes is used to attach all the routes to the API service.
func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		rate := limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		}
		cacheStore := persistence.NewInMemoryStore(time.Second)
		store := memory.NewStore()
		limiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
		api.Use(limiterMiddleware)

		coinGecko := services.NewCoinGeckoService()
		rateCtrl := controller.RateController{CoinGeckoService: coinGecko}
		api.GET(":coin", cache.CachePage(cacheStore, time.Minute*5, rateCtrl.GetCoin))
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
