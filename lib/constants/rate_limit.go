package constants

import "time"

const (
	// LoginRateLimit Login Rate Limit
	LoginRateLimit int64 = 10

	// LoginPeriod Login Periods
	LoginPeriod = 1 * time.Minute

	// BasicRateLimit Basic Rate Limit
	BasicRateLimit int64 = 200

	// BasicPeriod Basic Rate Period
	BasicPeriod = 1 * time.Minute
)
