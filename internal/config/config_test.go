package config

import (
	"testing"
)

func TestLoadFallsBackToDefaultsWhenEnvIsUnset(t *testing.T) {
	// t.Setenv cannot unset a variable, so clear each one explicitly and
	// rely on Setenv to restore it afterwards.
	for _, key := range []string{"APP_NAME", "APP_VERSION", "APP_ENV", "API_KEY"} {
		t.Setenv(key, "")
	}

	cfg := Load()

	if cfg.AppName != "go-k8s-app" {
		t.Errorf("AppName = %q, want %q", cfg.AppName, "go-k8s-app")
	}
	if cfg.AppVersion != "v1.0.0" {
		t.Errorf("AppVersion = %q, want %q", cfg.AppVersion, "v1.0.0")
	}
	if cfg.Environment != "development" {
		t.Errorf("Environment = %q, want %q", cfg.Environment, "development")
	}
	if cfg.APIKey != "default-secret" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "default-secret")
	}
}

func TestLoadPrefersEnvOverDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "from-env")
	t.Setenv("APP_VERSION", "v9.9.9")
	t.Setenv("APP_ENV", "staging")
	t.Setenv("API_KEY", "secret-from-env")

	cfg := Load()

	if cfg.AppName != "from-env" {
		t.Errorf("AppName = %q, want %q", cfg.AppName, "from-env")
	}
	if cfg.AppVersion != "v9.9.9" {
		t.Errorf("AppVersion = %q, want %q", cfg.AppVersion, "v9.9.9")
	}
	if cfg.Environment != "staging" {
		t.Errorf("Environment = %q, want %q", cfg.Environment, "staging")
	}
	if cfg.APIKey != "secret-from-env" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "secret-from-env")
	}
}

func TestGetEnvTreatsEmptyStringAsUnset(t *testing.T) {
	// Kubernetes delivers an unset ConfigMap key as an empty string, so an
	// empty value must not win over the fallback.
	t.Setenv("EMPTY_KEY", "")

	if got := getEnv("EMPTY_KEY", "fallback"); got != "fallback" {
		t.Errorf("getEnv(empty) = %q, want %q", got, "fallback")
	}
	if got := getEnv("MISSING_KEY", "fallback"); got != "fallback" {
		t.Errorf("getEnv(missing) = %q, want %q", got, "fallback")
	}
	t.Setenv("SET_KEY", "value")
	if got := getEnv("SET_KEY", "fallback"); got != "value" {
		t.Errorf("getEnv(set) = %q, want %q", got, "value")
	}
}
