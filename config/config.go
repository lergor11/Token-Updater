package config

import (
	"flag"
	"os"
)

type Config struct {
	VAULT_ROLE        string
	VAULT_PATH_ENV    string
	VAULT_ENV         string
	VAULT_AUTH_METHOD string
	API_URL           string
	API_METHOD        string
	API_CLIENT_ID     string
	API_CLIENT_SECRET string
	API_AUDINCE       string
	API_GRANT_TYPE    string
}

func MakeConfig() Config {
	var c Config

	VAULT_ROLE, found := os.LookupEnv("VAULT_ROLE")
	if !found {
		VAULT_ROLE = ""
	}
	flag.StringVar(&c.VAULT_ROLE, "vault.role", VAULT_ROLE,
		"Vault rw role")

	VAULT_AUTH_METHOD, found := os.LookupEnv("VAULT_AUTH_METHOD")
	if !found {
		VAULT_AUTH_METHOD = "kubernetes"
	}
	flag.StringVar(&c.VAULT_AUTH_METHOD, "vault.auth_metod", VAULT_AUTH_METHOD,
		"Vault auth metod")

	VAULT_PATH_ENV, found := os.LookupEnv("VAULT_PATH_ENV")
	if !found {
		VAULT_PATH_ENV = "example"
	}
	flag.StringVar(&c.VAULT_PATH_ENV, "vault.path_env", VAULT_PATH_ENV,
		"Vault path where save token")

	VAULT_ENV, found := os.LookupEnv("VAULT_ENV")
	if !found {
		VAULT_ENV = "token"
	}
	flag.StringVar(&c.VAULT_ENV, "vault.token", VAULT_ENV,
		"Where save token")

	API_URL, found := os.LookupEnv("API_URL")
	if !found {
		API_URL = "https://token"
	}
	flag.StringVar(&c.API_URL, "api.url", API_URL,
		"URL for token")

	API_METHOD, found := os.LookupEnv("API_METHOD")
	if !found {
		API_METHOD = "POST"
	}
	flag.StringVar(&c.API_METHOD, "api.method", API_METHOD,
		"Method api")

	API_CLIENT_ID, found := os.LookupEnv("API_CLIENT_ID")
	if !found {
		API_CLIENT_ID = ""
	}
	flag.StringVar(&c.API_CLIENT_ID, "api.client_id", API_CLIENT_ID,
		"API Client_ID")

	API_CLIENT_SECRET, found := os.LookupEnv("API_CLIENT_SECRET")
	if !found {
		API_CLIENT_SECRET = ""
	}
	flag.StringVar(&c.API_CLIENT_SECRET, "api.client_secret", API_CLIENT_SECRET,
		"API client_secret")

	API_AUDINCE, found := os.LookupEnv("API_AUDINCE")
	if !found {
		API_AUDINCE = ""
	}
	flag.StringVar(&c.API_AUDINCE, "api.audince", API_AUDINCE,
		"API Audience")

	API_GRANT_TYPE, found := os.LookupEnv("API_GRANT_TYPE")
	if !found {
		API_GRANT_TYPE = ""
	}
	flag.StringVar(&c.API_GRANT_TYPE, "api.grant_type", API_GRANT_TYPE,
		"API Grant_Type")

	flag.Parse()

	return c
}
