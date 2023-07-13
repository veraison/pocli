// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veraison/apiclient/management"
)

var (
	version = "0.0.1"

	configFile string
	config     = &Config{}

	rootCmd = &cobra.Command{
		Short:   "policy management client",
		Version: version,
	}

	service *management.Service
)

type Config struct {
	Host string
	Port int
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file")
	rootCmd.PersistentFlags().StringP("host", "H", "localhost",
		"the host running Veraison management service")
	rootCmd.PersistentFlags().IntP("port", "p", 8088,
		"the port on which Veraison management service is listening")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(activateCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deactivateCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(schemesCmd)
}

func init() {
	cobra.OnInitialize(func() {
		initConfig()
		initService()
	})
}

func initConfig() {
	v, err := readConfig(configFile)
	cobra.CheckErr(err)

	config.Host = v.GetString("host")
	config.Port = v.GetInt("port")
}

func readConfig(path string) (*viper.Viper, error) {
	v := viper.GetViper()
	if path != "" {
		v.SetConfigFile(path)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		userConfigDir, err := os.UserConfigDir()
		if err == nil {
			v.AddConfigPath(filepath.Join(userConfigDir, "pocli"))
		}
		v.AddConfigPath(wd)
		v.SetConfigType("yaml")
		v.SetConfigName("config")
	}

	v.SetEnvPrefix("pocli")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		err = nil
	}

	return v, err
}

func initService() {
	var err error

	serviceURI := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:   "/management/v1",
	}

	service, err = management.NewService(serviceURI.String())
	cobra.CheckErr(err)
}
