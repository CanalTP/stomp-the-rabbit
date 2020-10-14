package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Webstomp struct {
		Target      string
		Login       string
		Passcode    string
		Destination string
		Protocol    string
		SendTimeout int
		RecvTimeout int
	}
	AMQP struct {
		URL      string
		Exchange struct {
			Name string
		}
		ContentType string
	}
}

func (c *Config) ToString() string {
	var b strings.Builder

	fmt.Fprintln(&b, "Webstomp configuration")
	fmt.Fprintln(&b, "----------------------")
	fmt.Fprintln(&b, "Webstomp target is\t\t", c.Webstomp.Target)
	fmt.Fprintln(&b, "Webstomp login is\t\t", c.Webstomp.Login)
	fmt.Fprintln(&b, "Webstomp destination is\t\t", c.Webstomp.Destination)
	fmt.Fprintln(&b, "Webstomp protocol is\t\t", c.Webstomp.Protocol)
	fmt.Fprintln(&b, "Webstomp send timeout is\t", c.Webstomp.SendTimeout)
	fmt.Fprintln(&b, "Webstomp receive timeout is\t", c.Webstomp.RecvTimeout)
	fmt.Fprintln(&b, "\nAMQP configuration")
	fmt.Fprintln(&b, "------------------------")
	fmt.Fprintln(&b, "AMQP broker url is\t\t\t", c.AMQP.URL)
	fmt.Fprintln(&b, "Amqp exchange name is\t", c.AMQP.Exchange.Name)

	return b.String()
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("webstomp.sendTimeout", "0")
	viper.SetDefault("webstomp.recvTimeout", "0")
	viper.SetDefault("amqp.url", "amqp://guest:guest@localhost:5672//")
	viper.SetDefault("amqp.contentType", "application/json")

	// I need that to be able to unmarshal from env vars
	viper.BindEnv("webstomp.target")
	viper.BindEnv("webstomp.login")
	viper.BindEnv("webstomp.passcode")
	viper.BindEnv("webstomp.protocol")
	viper.BindEnv("webstomp.destination")
	viper.BindEnv("amqp.exchange.name")
}

func config() (*Config, error) {
	// handle error while reading configuration file
	// if file not found -> no error since this file is optionnal
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	C := new(Config)
	err := viper.Unmarshal(C)
	if err != nil {
		return nil, err
	}

	return C, nil
}
