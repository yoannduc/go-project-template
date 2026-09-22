package env

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/yoannduc/go-project-template/pkg/singleton"
)

func TestParseEnv(t *testing.T) {
	testCases := []struct {
		env string
		out Env
		err error
	}{
		{"", EnvProd, errInvalidEnv},
		{"string", EnvProd, errInvalidEnv},
		{"STRING", EnvProd, errInvalidEnv},
		{"PROD", EnvProd, nil},
		{"prod", EnvProd, nil},
		{"pRoD", EnvProd, nil},
		{"PRODUCTION", EnvProd, nil},
		{"DEV", EnvDev, nil},
		{"dev", EnvDev, nil},
		{"DEVELOPMENT", EnvDev, nil},
		{"TEST", EnvTest, nil},
		{"test", EnvTest, nil},
		{stringEnvDev, EnvDev, nil},
		{stringEnvDevelopment, EnvDev, nil},
		{stringEnvProd, EnvProd, nil},
		{stringEnvProduction, EnvProd, nil},
		{stringEnvTest, EnvTest, nil},
	}

	for _, test := range testCases {
		t.Run(test.env, func(t *testing.T) {
			v, err := parseEnv(test.env)
			if test.out != v {
				t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.env, test.out, v)
			}
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error parsing env. Expected "%v", got "%v"`, test.err, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Run("existing value", func(t *testing.T) {
		testCases := map[string]struct {
			env Env
			out Env
		}{
			"-100":        {Env(-100), Env(-100)},
			"-1":          {Env(-1), Env(-1)},
			"EnvTest - 1": {EnvTest - 1, EnvTest - 1},
			"EnvTest":     {EnvTest, EnvTest},
			"EnvDev":      {EnvDev, EnvDev},
			"EnvProd":     {EnvProd, EnvProd},
			"EnvProd + 1": {EnvProd + 1, EnvProd + 1},
			"1":           {Env(1), Env(1)},
			"100":         {Env(100), Env(100)},
		}

		for name, test := range testCases {
			// Save original env
			originalEnv := singl

			t.Run(name, func(t *testing.T) {
				singl = &singleton.Singleton[Env]{Instance: test.env}

				v := Get()
				if test.out != v {
					t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.env, test.out, v)
				}
			})

			// Revert original env
			singl = originalEnv
		}
	})

	t.Run("parse env from env", func(t *testing.T) {
		testCases := []struct {
			env string
			out Env
		}{
			{"", EnvProd},
			{"string", EnvProd},
			{"STRING", EnvProd},
			{"PROD", EnvProd},
			{"PRODUCTION", EnvProd},
			{"DEV", EnvDev},
			{"dev", EnvDev},
			{"dEv", EnvDev},
			{"DEVELOPMENT", EnvDev},
			{"development", EnvDev},
			{"TEST", EnvTest},
			{"test", EnvTest},
		}

		for _, test := range testCases {
			// Save original env
			originalEnv := os.Getenv(varName)

			t.Run(test.env, func(t *testing.T) {
				// Set env to input for test & reset once & reset singleton
				_ = os.Setenv(varName, test.env)
				once = sync.Once{}
				singl = nil

				v := Get()
				if test.out != v {
					t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.env, test.out, v)
				}
			})

			// Revert original env & reset once & reset singleton
			_ = os.Setenv(varName, originalEnv)
			once = sync.Once{}
			singl = nil
		}
	})

	t.Run("env var read only once", func(t *testing.T) {
		// Save original env
		originalEnv := os.Getenv(varName)

		// Set env to input for test & reset once & reset singleton
		_ = os.Setenv(varName, stringEnvDev)
		once = sync.Once{}
		singl = nil

		v := Get()
		if v != EnvDev {
			t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, stringEnvDev, EnvDev, v)
		}

		_ = os.Setenv(varName, stringEnvProd)

		v = Get()
		if v != EnvDev {
			t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, stringEnvProd, EnvDev, v)
		}

		// Revert original env & reset once & reset singleton
		_ = os.Setenv(varName, originalEnv)
		once = sync.Once{}
		singl = nil
	})
}

func TestIsDev(t *testing.T) {
	testCases := map[string]struct {
		env Env
		out bool
	}{
		"-100":        {Env(-100), true},
		"-1":          {Env(-1), true},
		"EnvTest - 1": {EnvTest - 1, true},
		"EnvTest":     {EnvTest, true},
		"EnvDev":      {EnvDev, true},
		"EnvProd":     {EnvProd, false},
		"EnvProd + 1": {EnvProd + 1, false},
		"1":           {Env(1), false},
		"100":         {Env(100), false},
	}

	for name, test := range testCases {
		// Save original env
		originalEnv := singl

		t.Run(name, func(t *testing.T) {
			singl = &singleton.Singleton[Env]{Instance: test.env}

			v := IsDev()
			if test.out != v {
				t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.env, test.out, v)
			}
		})

		// Revert original env
		singl = originalEnv
	}
}
