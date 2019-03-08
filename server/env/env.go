// Package env defines the default environment variables used within the app.
package env

import (
	"os"
)

func getEnv(key, def string) string {
	k := os.Getenv(key)
	if k == "" {
		return def
	}
	return k
}

var isProd bool
var stage string

func init() {
	stage = getEnv("GO_ENV", "development")
	isProd = stage == "production"
}

//Stage returns GO_ENV defaults to development
func Stage() string { return stage }

//IsProd returns true if GO_ENV=production
func IsProd() bool { return isProd }
