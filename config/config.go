package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	UnmarshalKey(string, interface{}, ...viper.DecoderConfigOption) error
}

var defaultConfig *viper.Viper

// Config returns a default config providers
func Config() Provider {
	return defaultConfig
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig("UPPER_gutenberg-ingester")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("download_type", ".epub3.images")
	v.SetDefault("update_previously_downloaded", false)
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()

	setDefaults(v)

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath(".")
	v.AddConfigPath("/etc/gutenberg-ingester/")

	v.ReadInConfig()

	v.SetEnvPrefix(appName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// workaround because viper does not treat env vars the same as other config
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		v.Set(key, val)
	}

	return v
}
