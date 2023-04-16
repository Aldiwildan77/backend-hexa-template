package config

import "time"

type Middleware struct {
	RateLimit RateLimit
	Timeout   Timeout
	CORS      CORS
}

type RateLimit struct {
	Rate         int    `env:"RATE_LIMIT_RATE" envDefault:"10"`
	Burst        int    `env:"RATE_LIMIT_BURST" envDefault:"10"`
	ExpireTime   int    `env:"RATE_LIMIT_EXPIRE" envDefault:"60"`
	ExpireFormat string `env:"RATE_LIMIT_EXPIRE_FORMAT" envDefault:"s"`
}

// GetExpireTime return time in expire time format
func (rl RateLimit) GetExpireTime() time.Duration {
	return time.Duration(rl.ExpireTime) * rl.GetExpireFormat()
}

// GetExpireFormat return time.Duration from input string
func (rl RateLimit) GetExpireFormat() time.Duration {
	switch rl.ExpireFormat {
	case "s", "sec", "second":
		return time.Second
	case "m", "min", "minute":
		return time.Minute
	case "h", "hour":
		return time.Hour
	case "d", "day":
		return time.Hour * 24
	default:
		return time.Second
	}
}

type CORS struct {
	AllowOrigins []string `env:"CORS_ALLOW_ORIGINS" envDefault:"*" envSeparator:","`
	AllowMethods []string `env:"CORS_ALLOW_METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS" envSeparator:","`
}

type Timeout struct {
	Duration int `env:"HTTP_TIMEOUT_DURATION" envDefault:"10"`
}
