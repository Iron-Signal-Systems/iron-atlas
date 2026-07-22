package app

import (
	"testing"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
)

func clearAuthenticationEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("IRON_ATLAS_AUTHENTICATION_MODE", "")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "")
}

func TestConfigDefaultsToMemoryDevelopmentMode(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", "")
	t.Setenv("IRON_ATLAS_DATABASE_URL", "")
	clearAuthenticationEnvironment(t)

	cfg, err := ConfigFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ChangeStore != ChangeStoreMemory {
		t.Fatalf("expected memory store, got %s", cfg.ChangeStore)
	}
	if cfg.AuthenticationMode != authentication.ModeDevelopment {
		t.Fatalf(
			"memory demonstration mode = %q, want development",
			cfg.AuthenticationMode,
		)
	}
}

func TestPostgreSQLModeRequiresURL(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", ChangeStorePostgreSQL)
	t.Setenv("IRON_ATLAS_DATABASE_URL", "")
	clearAuthenticationEnvironment(t)

	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("expected missing database URL error")
	}
}

func TestPostgreSQLModeDefaultsToProductionAuthentication(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", ChangeStorePostgreSQL)
	t.Setenv("IRON_ATLAS_DATABASE_URL", "postgres://example")
	clearAuthenticationEnvironment(t)

	cfg, err := ConfigFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationMode != authentication.ModeProduction {
		t.Fatalf(
			"persistent mode = %q, want production",
			cfg.AuthenticationMode,
		)
	}
}

func TestDevelopmentModeCanBeExplicitlyEnabledForControlledTesting(
	t *testing.T,
) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", ChangeStorePostgreSQL)
	t.Setenv("IRON_ATLAS_DATABASE_URL", "postgres://example")
	t.Setenv("IRON_ATLAS_AUTHENTICATION_MODE", "development")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "")

	cfg, err := ConfigFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationMode != authentication.ModeDevelopment {
		t.Fatalf(
			"explicit mode = %q, want development",
			cfg.AuthenticationMode,
		)
	}
}

func TestInvalidAuthenticationModeIsRejected(t *testing.T) {
	t.Setenv("IRON_ATLAS_AUTHENTICATION_MODE", "sometimes")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "")

	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("expected invalid authentication mode error")
	}
}

func TestLegacyDevelopmentIdentitySettingIsRejected(t *testing.T) {
	t.Setenv("IRON_ATLAS_AUTHENTICATION_MODE", "")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "true")

	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("legacy boolean identity setting must fail closed")
	}
}

func TestPoolLimitsAreValidatedAsIntegers(t *testing.T) {
	t.Setenv("IRON_ATLAS_DATABASE_MAX_CONNECTIONS", "not-an-integer")
	clearAuthenticationEnvironment(t)

	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("expected invalid maximum connection error")
	}
}
