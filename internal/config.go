package internal

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ConfigStore defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type ConfigStore interface {
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
func Config() ConfigStore {
	return readViperConfig("GUTENBERG_INGESTER")
}

func DevConfig() ConfigStore {
	return readViperDevConfig("GUTENBERG_INGESTER")
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) ConfigStore {
	return readViperConfig(appName)
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("mode", "production")
	v.SetDefault("database_location", "/var/gutenberg-ingester/db.json")
	v.SetDefault("temporary_directory", "/tmp")
	v.SetDefault("download_location", "./")
	v.SetDefault("download_type", []string{".epub3.images", ".epub", ".txt"})
	v.SetDefault("download_type_precedence", "sequential")
	v.SetDefault("download_delay", 2)
	v.SetDefault("download_languages", []string{"english"})
	v.SetDefault("update_previously_downloaded", false)
	v.SetDefault("gutenberg_feed_url", "https://www.gutenberg.org/")
	v.SetDefault("gutenberg_mirror_url", "https://www.gutenberg.org/")
	v.SetDefault("full_sync_frequency", 7)
	v.SetDefault("partial_sync_frequency", 12)
	v.SetDefault("epub_use_proper_extension", false)
}

func setDevOverideDefaults(v *viper.Viper) {
	v.SetDefault("mode", "development")
	v.SetDefault("database_location", "/tmp/gutenberg-ingester-db.json")
	v.SetDefault("download_location", "/tmp/")
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

func readViperDevConfig(appName string) *viper.Viper {
	v := viper.New()

	setDefaults(v)
	setDevOverideDefaults(v)

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
