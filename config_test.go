package kocto_test

import (
	"os"
	"testing"

	"github.com/Kamae-RT/kocto"
	"github.com/matryer/is"
)

func TestConfig(t *testing.T) {
	is := is.NewRelaxed(t)

	os.Setenv("ENV", "production")
	os.Setenv("PORT", "4123")
	os.Setenv("LOG_NAME", "test log")
	os.Setenv("DATABASE_NAME", "test database")

	cfg, err := kocto.LoadConfig()

	is.NoErr(err)
	is.Equal(cfg.Env, kocto.Production)
	is.Equal(cfg.Port, "4123")
	is.Equal(cfg.Log.Name, "test log")
	is.Equal(cfg.DB.Name, "test database")
}

func TestExtendedConfig(t *testing.T) {
	is := is.NewRelaxed(t)

	os.Setenv("ENV", "production")
	os.Setenv("PORT", "4123")
	os.Setenv("LOG_NAME", "test log")
	os.Setenv("DATABASE_NAME", "test database")
	os.Setenv("MY_INNER_PROP", "inner prop")
	os.Setenv("MY_NESTED_INNER_PROP", "nested inner prop")

	var cfg struct {
		kocto.Config
		InnerProp string `env:"MY_INNER_PROP"`
		Nested    struct {
			NesterInner string `env:"MY_NESTED_INNER_PROP"`
		}
	}

	err := kocto.LoadInConfig(&cfg)

	is.NoErr(err)
	is.Equal(cfg.Env, kocto.Production)
	is.Equal(cfg.Port, "4123")
	is.Equal(cfg.Log.Name, "test log")
	is.Equal(cfg.DB.Name, "test database")
	is.Equal(cfg.InnerProp, "inner prop")
	is.Equal(cfg.Nested.NesterInner, "nested inner prop")
}
