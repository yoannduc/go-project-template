package env

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/yoannduc/go-project-template/pkg/logger"
	"github.com/yoannduc/go-project-template/pkg/singleton"
)

var (
	errInvalidEnv = errors.New("invalid environment value")
)

type Env int

const (
	EnvTest Env = -2
	EnvDev  Env = -1
	EnvProd Env = 0
)

const (
	stringEnvTest        = "TEST"
	stringEnvDev         = "DEV"
	stringEnvDevelopment = "DEVELOPMENT"
	stringEnvProd        = "PROD"
	stringEnvProduction  = "PRODUCTION"
)

// parseEnv returns the env level based on env string.
// It returns production level and an error if string did not
// match any env handled. It is case-insensitive.
func parseEnv(env string) (Env, error) {
	switch strings.ToUpper(env) {
	case stringEnvTest:
		return EnvTest, nil
	case stringEnvDev, stringEnvDevelopment:
		return EnvDev, nil
	case stringEnvProd, stringEnvProduction:
		return EnvProd, nil
	default:
		return EnvProd, errInvalidEnv
	}
}

const (
	// varName is the env var name used to determine environment.
	// It is a const for tests.
	varName = "ENV"
)

var (
	once  sync.Once
	singl *singleton.Singleton[Env]
)

// Get returns the env level set by env variable [varName].
// It uses singleton design pattern to read from env only once.
// If no env value is set in env, it warns and defaults to production level.
func Get() Env {
	if singl == nil {
		once.Do(func() {
			env := os.Getenv(varName)
			v, err := parseEnv(env)
			if err != nil {
				logger.Get().LogAttrs(context.Background(), logger.LevelWarn, "defaulting to "+stringEnvProd, slog.String("env", env), slog.Any("err", err))
			}

			singl = &singleton.Singleton[Env]{Instance: v}
		})
	}

	return singl.Instance
}

// IsDev returns whether the env is dev/test mode. It uses [Get].
func IsDev() bool {
	return Get() < EnvProd
}
