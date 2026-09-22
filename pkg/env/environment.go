package env

import (
	"errors"
	"os"
	"strings"
)

var (
	errInvalidEnv = errors.New("invalid environment value")
)

// An Env represents the environment level of the service.
type Env int

// Names for common environment levels.
const (
	EnvTest Env = -2
	EnvDev  Env = -1
	EnvProd Env = 0
)

// IsDev returns whether the env is dev/test mode.
func (env Env) IsDev() bool { return env < EnvProd }

// IsTest returns whether the env is test mode only.
func (env Env) IsTest() bool { return env < EnvDev }

// Strings for handled environment.
const (
	stringEnvTest        = "TEST"
	stringEnvDev         = "DEV"
	stringEnvDevelopment = "DEVELOPMENT"
	stringEnvProd        = "PROD"
	stringEnvProduction  = "PRODUCTION"
	stringEnvUnsupported = "UNSUPPORTED"
)

// stringifyEnv returns the strign representation of env level.
// It returns stringEnvUnsupported and an error if env did not
// match any env handled.
func stringifyEnv(env Env) (string, error) {
	switch env {
	case EnvTest:
		return stringEnvTest, nil
	case EnvDev:
		return stringEnvDev, nil
	case EnvProd:
		return stringEnvProd, nil
	default:
		return stringEnvUnsupported, errInvalidEnv
	}
}

// String returns the strign representation of env level.
func (env Env) String() string {
	s, _ := stringifyEnv(env)
	return s
}

// parseEnv returns the env level based on env string.
// It returns production level and an error if string did not
// match any env handled or is empty. It is case-insensitive.
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
	varName = "ENV"
)

// Get returns the env level set by env variable defined in varName.
// If no env value is set in env, it defaults to production level.
func Get() Env {
	v, _ := parseEnv(os.Getenv(varName))

	return v
}
