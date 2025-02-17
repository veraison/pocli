// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/veraison/apiclient/auth"
	"github.com/veraison/apiclient/management"
)

var (
	version = "0.0.1"

	configFile string
	config     = &Config{}

	authMethod = auth.MethodPassthrough

	rootCmd = &cobra.Command{
		Short:   "policy management client",
		Version: version,
	}

	service *management.Service
)

type Config struct {
	Host string
	Port int

	Auth auth.IAuthenticator

	Insecure bool
	UseTLS bool
	CertPaths []string
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file")
	rootCmd.PersistentFlags().StringP("host", "H", "localhost",
		"the host running Veraison management service")
	rootCmd.PersistentFlags().IntP("port", "p", 10443,
		"the port on which Veraison management service is listening")
	rootCmd.PersistentFlags().VarP(&authMethod, "auth", "a",
		`authentication method, must be one of "none"/"passthrough", "basic", "oauth2"`)
	rootCmd.PersistentFlags().StringP("client-id", "C", "", "OAuth2 client ID")
	rootCmd.PersistentFlags().StringP("client-secret", "S", "", "OAuth2 client secret")
	rootCmd.PersistentFlags().StringP("token-url", "T", "", "token URL of the OAuth2 service")
	rootCmd.PersistentFlags().StringP("username", "U", "", "service username")
	rootCmd.PersistentFlags().StringP("password", "P", "", "service password")
	rootCmd.PersistentFlags().BoolP(
		"insecure", "i", false, "Use HTTPS but do not check certs (implies -s/--tls)",
	)
	rootCmd.PersistentFlags().BoolP(
		"tls", "s", false, "Use HTTPS",
	)
	rootCmd.PersistentFlags().StringArrayP(
		"ca-cert", "E", nil, "path to a CA cert that will be used in addition to system certs; may be specified multiple times",
	)

	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		cfgName := strings.ReplaceAll(flag.Name, "-", "_")
		err := viper.BindPFlag(cfgName, flag)
		cobra.CheckErr(err)
	})

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
	config.Insecure = v.GetBool("insecure")
	config.UseTLS = v.GetBool("tls")
	config.CertPaths = v.GetStringSlice("ca_cert")

	err = authMethod.Set(v.GetString("auth"))
	cobra.CheckErr(err)

	switch authMethod {
	case auth.MethodPassthrough:
		config.Auth = &auth.NullAuthenticator{}
	case auth.MethodBasic:
		config.Auth = &auth.BasicAuthenticator{}
		err = config.Auth.Configure(map[string]interface{}{
			"username": v.GetString("username"),
			"password": v.GetString("password"),
		})
		cobra.CheckErr(err)
	case auth.MethodOauth2:
		config.Auth = &auth.Oauth2Authenticator{}
		err = config.Auth.Configure(map[string]interface{}{
			"client_id":     v.GetString("client_id"),
			"client_secret": v.GetString("client_secret"),
			"token_url":     v.GetString("token_url"),
			"username":      v.GetString("username"),
			"password":      v.GetString("password"),
			"ca_certs":      v.GetStringSlice("ca_cert"),
		})
		cobra.CheckErr(err)
	default:
		// Should never get here as authMethod value is set via
		// Method.Set(), which ensures that it's one of the above.
		panic(fmt.Sprintf("unknown auth method: %q", authMethod))
	}

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

	scheme := "http"
	if config.Insecure || config.UseTLS {
		scheme = "https"
	}

	serviceURI := url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:   "/management/v1",
	}

	if config.Insecure {
		service, err = management.NewInsecureTLSService(serviceURI.String(), config.Auth)
	} else if config.UseTLS {
		service, err = management.NewTLSService(
			serviceURI.String(),
			config.Auth,
			config.CertPaths,
		)
	} else {
		service, err = management.NewService(serviceURI.String(), config.Auth)
	}
	cobra.CheckErr(err)
}
