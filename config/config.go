package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
)

type Domain struct {
	Main []string `yaml:"main"`
}

type Http struct {
	Domains Domain `yaml:"domains"`
}

type CircuitBreaker struct {
	FailureThreshold int `yaml:"failureThreshold"`
	SuccessThreshold int `yaml:"successThreshold"`
	Timeout          int `yaml:"timeout"`
}

type EntryPoints struct {
	Address string `yaml:"address"`
	Http    Http   `yaml:"http"`
}

type config struct {
	HttpClient     HttpClient     `yaml:"ServiceHttp"`
	GRPCClient     GRPCClient     `yaml:"ServiceGRPC"`
	EntryPoints    EntryPoints    `yaml:"EntryPoints"`
	App            App            `yaml:"App"`
	AppPort        AppPort        `yaml:"AppPort"`
	CircuitBreaker CircuitBreaker `yaml:"CircuitBreaker"`
}

type HttpClient struct {
	Port string `yaml:"PORT"`
	Host string `yaml:"HOST"`
}

type GRPCClient struct {
	Port string `yaml:"PORT"`
	Host string `yaml:"HOST"`
}

type AppPort struct {
	Port string `yaml:"PORT"`
}

type App struct {
	LogLevel bool `yaml:"LOG_LEVEL"`
}

var cfg config

type ConfiguredApp struct {
	Config config
	PORT   string
}

func InitializeConfig() *ConfiguredApp {
	viper.SetEnvPrefix("GATEWAY")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.MergeInConfig()
	if err != nil {
		fmt.Println("Error in reading Config")
		panic(err)
	}

	err = viper.Unmarshal(&cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
	})
	if err != nil {
		fmt.Println("Error in unmarshalling Config")
		panic(err)
	}

	fmt.Printf("\n loaded Config: %#v \n", cfg)

	return &ConfiguredApp{
		Config: cfg,
		PORT:   cfg.AppPort.Port,
	}
}
