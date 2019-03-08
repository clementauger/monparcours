package app

import (
	"fmt"
	"math/rand"

	"github.com/clementauger/monparcours/server/config"
)

//Environment defines the execution options.
type Environment struct {
	// CanonicalHost of the website
	CanonicalHost string `yaml:"canonicalhost"`
	// Port listened by the server.
	Port int `yaml:"port"`
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
	GlobalRateLimit *RateLimit `yaml:"globalratelimit"`
	//LoginRateLimit applies to the admin login.
	LoginRateLimit *RateLimit `yaml:"loginratelimit"`
	//GeoCoderCacheSize defines the length of the LRU cache of osm requests.
	GeoCoderCacheSize int `yaml:"geocodercacheSize"`
}

//RateLimit configures http rate limiter.
type RateLimit struct {
	Size  int `yaml:"size"`
	RPM   int `yaml:"rpm"`
	Burst int `yaml:"burst"`
}

//GetEnvironment loads the config and applies default values.
func GetEnvironment(filename, environment string) (*Environment, error) {
	conf := make(map[string]*Environment)
	err := config.ReadConfig(conf, filename)
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

	if env.GlobalRateLimit == nil {
		env.GlobalRateLimit = &RateLimit{
			Burst: 5,
			Size:  65536,
			RPM:   20,
		}
	}
	if env.LoginRateLimit == nil {
		env.LoginRateLimit = &RateLimit{
			Burst: 1,
			Size:  65536,
			RPM:   10,
		}
	}

	return env, nil
}

//GetApp returns an http application to serve.
func GetApp(env *Environment) (HTTPApp, error) {
	app := HTTPApp{Env: *env}
	return app, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
