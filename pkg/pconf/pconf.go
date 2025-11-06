package pconf

import (
	"io"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// 配置优先级
	// set > flag > env > file
	//viper.AddConfigPath("etc")
	//viper.AddConfigPath("/etc")
	//viper.SetEnvPrefix()
	viper.AutomaticEnv()
}

// SetConfigName sets name for the config file
func SetConfigName(in string) { viper.SetConfigName(in) }

// SetConfigType sets the type of the configuration
func SetConfigType(in string) { viper.SetConfigType(in) }

// AddConfigPath adds a path for Viper to search for the config file
func AddConfigPath(in string) { viper.AddConfigPath(in) }

// SetConfigFile explicitly defines the path, name and extension of the config file
func SetConfigFile(in string) { viper.SetConfigFile(in) }

// ReadInConfig will auto discover and load the configuration file from disk
func ReadInConfig() error { return viper.ReadInConfig() }

func ReadConfig(in io.Reader) error  { return viper.ReadConfig(in) }
func MergeConfig(in io.Reader) error { return viper.MergeConfig(in) }

// WriteConfig writes the current configuration to a file
func WriteConfig() error { return viper.WriteConfig() }

// SetWatchChange will call back when config changed
func SetWatchChange(fn func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fn()
	})
}

// Set sets the value for the key
func Set(key string, value interface{}) { viper.Set(key, value) }

// SetDefault sets the default value for this key
func SetDefault(key string, value interface{}) { viper.SetDefault(key, value) }

func GetBool(key string) bool     { return viper.GetBool(key) }
func GetInt(key string) int       { return viper.GetInt(key) }
func GetString(key string) string { return viper.GetString(key) }

// AutomaticEnv 自动加载环境变量
func AutomaticEnv() { viper.AutomaticEnv() }

// SetEnvPrefix 设置环境变量前缀：
// SetEnvPrefix("pconf")
// will auto load env: PCONF_xxx
func SetEnvPrefix(in string) { viper.SetEnvPrefix(in) }

// BindPFlag 自动绑定命令行参数
// serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
// viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
func BindPFlag(key string, flag *pflag.Flag) error { return viper.BindPFlag(key, flag) }

// Unmarshal 反序列化
// tag default: `mapstructure:"item_name"`
// key like: app.config
func Unmarshal(rawVal interface{}) error {
	opt := func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	}
	return viper.Unmarshal(rawVal, opt)
}
func UnmarshalKey(key string, rawVal interface{}) error {
	opt := func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	}
	return viper.UnmarshalKey(key, rawVal, opt)
}
