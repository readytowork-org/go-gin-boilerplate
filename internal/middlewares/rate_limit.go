package middlewares

import (
	"boilerplate-api/internal/json_response"
	"net/http"
	"strconv"
	"time"

	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// TODO :: refactor

// Global store
// using in-memory store with goroutine which clears expired keys.
var store = memory.NewStore()

type RateLimitOption struct {
	period time.Duration
	limit  int64
}

type Option func(*RateLimitOption)

type RateLimitMiddleware struct {
	logger config.Logger
}

func NewRateLimitMiddleware(logger config.Logger) RateLimitMiddleware {
	return RateLimitMiddleware{
		logger: logger,
	}
}

func (rl RateLimitMiddleware) HandleRateLimit(limit int64, period time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() // Gets cient IP Address

		rl.logger.Info("Setting up rate limit middleware...")

		// Limit # of API Calls
		// Period in a given time frame
		// setting default values
		opt := RateLimitOption{
			period: period,
			limit:  limit,
		}

		rate := limiter.Rate{
			Limit:  opt.limit,
			Period: opt.period,
		}

		instance := limiter.New(store, rate)

		context, err := instance.Get(c, c.FullPath()+"&&"+key)

		if err != nil {
			rl.logger.Panic(err.Error())
		}

		c.Set(constants.RateLimit, instance)

		// Setting custom headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Limit exceeded
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, json_response.Error[string]{
				Error:   "Too many request",
				Message: "Rate limit has exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func WithOptions(period time.Duration, limit int64) Option {
	return func(o *RateLimitOption) {
		o.period = period
		o.limit = limit
	}
}
