package app

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/clementauger/coders"
	"github.com/clementauger/httpr"
)

//Environment defines the execution options.
type Environment struct {
	// CanonicalHost of the website
	CanonicalHost string `yaml:"canonicalhost"`
	// Port listened by the server.
	Port int `yaml:"port"`
	// ReadTimeut of http incoming request.
	ReadTimeout *time.Duration `yaml:"readtimeout"`
	// WriteTimeut of http outgoing request.
	WriteTimeut *time.Duration `yaml:"writetimeout"`
	//CsrfKey to protext the tokens (default 32 bytes random string)
	CsrfKey string `yaml:"csrf"`
	//AdminKey to generate an admin password (default 32 bytes random string)
	AdminKey string `yaml:"adminkey"`
	//AdminSalt to generate an admin password (default 32 bytes random string)
	AdminSalt string `yaml:"adminsalt"`
	// SchemaName        string     `yaml:"schema"`
	//CaptchaSolution for development and testing purposes.
	CaptchaSolution string `yaml:"captchasolution"`
	//Statik defiens if the application loads its assets from the regular file system ot the embedded static assets.
	Statik bool `yaml:"statik"`
	//PwdSalt makes user password stronger.
	PwdSalt string `yaml:"pwdsalt"`
	//Host to serve, it can a regexp as defiend by the gorilla framework.
	Host string `yaml:"host"`
	//GlobalRateLimit applies on all web (default 32 bytes random string)site and all requests.
	GlobalRateLimit *httpr.RateLimit `yaml:"globalratelimit"`
	//LoginRateLimit applies to the admin login.
	LoginRateLimit *httpr.RateLimit `yaml:"loginratelimit"`
	//GeoCoderCacheSize defines the length of the LRU cache of osm requests.
	GeoCoderCacheSize int `yaml:"geocodercachesize"`
}

//GetEnvironment loads the config and applies default values.
func GetEnvironment(filename, environment string) (*Environment, error) {
	conf := make(map[string]*Environment)
	err := coders.Decode(conf, filename)
	if err != nil {
		return nil, err
	}

	env := conf[environment]
	if env == nil {
		return nil, fmt.Errorf("environment %q does not have configuration", environment)
	}

	if env.CsrfKey == "" {
		env.CsrfKey = randStringRunes(32)
	}
	if env.AdminKey == "" {
		env.AdminKey = randStringRunes(32)
	}
	if env.AdminSalt == "" {
		env.AdminSalt = randStringRunes(32)
	}
	if env.PwdSalt == "" {
		env.PwdSalt = randStringRunes(32)
	}
	if env.GeoCoderCacheSize < 1 {
		env.GeoCoderCacheSize = 1024
	}

	if env.GlobalRateLimit == nil {
		env.GlobalRateLimit = &httpr.RateLimit{
			Burst: 5,
			Size:  65536,
			RPM:   20,
		}
	}
	if env.LoginRateLimit == nil {
		env.LoginRateLimit = &httpr.RateLimit{
			Burst: 1,
			Size:  65536,
			RPM:   10,
		}
	}
	if env.ReadTimeout == nil {
		y := time.Second * 15
		env.ReadTimeout = &y
	}
	if env.WriteTimeut == nil {
		y := time.Second * 15
		env.WriteTimeut = &y
	}

	return env, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
