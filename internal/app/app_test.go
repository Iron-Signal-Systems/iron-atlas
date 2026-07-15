package app

import "testing"

func TestConfigDefaultsToMemoryDevelopmentMode(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", "")
	t.Setenv("IRON_ATLAS_DATABASE_URL", "")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "")
	cfg, err := ConfigFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ChangeStore != ChangeStoreMemory {
		t.Fatalf("expected memory store, got %s", cfg.ChangeStore)
	}
	if !cfg.DevelopmentIdentity {
		t.Fatal("memory demonstration mode should default development identity on")
	}
}

func TestPostgreSQLModeRequiresURL(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", ChangeStorePostgreSQL)
	t.Setenv("IRON_ATLAS_DATABASE_URL", "")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "")
	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("expected missing database URL error")
	}
}

func TestPostgreSQLModeDefaultsDevelopmentIdentityOff(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", ChangeStorePostgreSQL)
	t.Setenv("IRON_ATLAS_DATABASE_URL", "postgres://example")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "")
	cfg, err := ConfigFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DevelopmentIdentity {
		t.Fatal("persistent mode must not silently trust development identity headers")
	}
}

func TestDevelopmentIdentityCanBeExplicitlyEnabledForControlledTesting(t *testing.T) {
	t.Setenv("IRON_ATLAS_CHANGE_STORE", ChangeStorePostgreSQL)
	t.Setenv("IRON_ATLAS_DATABASE_URL", "postgres://example")
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "true")
	cfg, err := ConfigFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.DevelopmentIdentity {
		t.Fatal("explicit development identity setting was not honored")
	}
}

func TestInvalidDevelopmentIdentitySettingIsRejected(t *testing.T) {
	t.Setenv("IRON_ATLAS_DEV_IDENTITY", "sometimes")
	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("expected invalid development identity setting error")
	}
}

func TestPoolLimitsAreValidatedAsIntegers(t *testing.T) {
	t.Setenv("IRON_ATLAS_DATABASE_MAX_CONNECTIONS", "not-an-integer")
	if _, err := ConfigFromEnvironment(); err == nil {
		t.Fatal("expected invalid maximum connection error")
	}
}
