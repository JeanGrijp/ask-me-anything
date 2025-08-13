package configAuth

import (
	"os"
	"time"
)

var (
	JWTSecret = os.Getenv("JWT_SECRET")

	JWTTokenExpires = func() time.Duration {
		val := os.Getenv("JWT_EXPIRES")
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
		return time.Hour * 24
	}()
)
