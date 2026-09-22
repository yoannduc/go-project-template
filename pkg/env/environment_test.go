package env

import (
	"errors"
	"testing"
)

func TestIsDev(t *testing.T) {
	testCases := map[string]struct {
		in  Env
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
		t.Run(name, func(t *testing.T) {
			v := test.in.IsDev()
			if test.out != v {
				t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestIsTest(t *testing.T) {
	testCases := map[string]struct {
		in  Env
		out bool
	}{
		"-100":        {Env(-100), true},
		"EnvTest - 1": {EnvTest - 1, true},
		"EnvTest":     {EnvTest, true},
		"-1":          {Env(-1), false},
		"EnvDev":      {EnvDev, false},
		"EnvProd":     {EnvProd, false},
		"EnvProd + 1": {EnvProd + 1, false},
		"1":           {Env(1), false},
		"100":         {Env(100), false},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			v := test.in.IsTest()
			if test.out != v {
				t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

func TestStringifyEnv(t *testing.T) {
	testCases := []struct {
		in  Env
		out string
		err error
	}{
		{Env(1000), stringEnvUnsupported, errInvalidEnv},
		{EnvDev, stringEnvDev, nil},
		{EnvProd, stringEnvProd, nil},
		{EnvTest, stringEnvTest, nil},
	}

	for _, test := range testCases {
		t.Run(test.in.String(), func(t *testing.T) {
			v, err := stringifyEnv(test.in)
			if test.out != v {
				t.Fatalf(`wrong result for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
			if err != nil && !errors.Is(err, test.err) {
				t.Fatalf(`unexpected error. Expected "%v", got "%v"`, test.err, err)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		in  Env
		out string
	}{
		{Env(1000), stringEnvUnsupported},
		{EnvDev, stringEnvDev},
		{EnvProd, stringEnvProd},
		{EnvTest, stringEnvTest},
	}

	for _, test := range testCases {
		t.Run(test.in.String(), func(t *testing.T) {
			v := test.in.String()
			if test.out != v {
				t.Fatalf(`wrong result for input "%v". Expected %v, got %v`, test.in, test.out, v)
			}
		})
	}
}

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
		t.Run(test.env, func(t *testing.T) {
			t.Setenv(varName, test.env)
			v := Get()
			if test.out != v {
				t.Fatalf(`wrong result for env "%v". Expected %v, got %v`, test.env, test.out, v)
			}
		})
	}
}
