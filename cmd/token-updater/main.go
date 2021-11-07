package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"io/ioutil"
	"net/http"
	"os"
	"token-updater/config"
)

var DefaultConfig config.Config

func init() {
	DefaultConfig = config.MakeConfig()
}

func Run(cfg config.Config) {
	token := getToken(cfg)

	err := putSecretWithKubernetesAuth(cfg, token)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to initialize Vault client: %w", err))
	}

}

func putSecretWithKubernetesAuth(cfg config.Config, secretEnv string) (err error) {
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to initialize Vault client: %w", err))
	}

	// Read the service-account token from the path where the token's Kubernetes Secret is mounted.
	// By default, Kubernetes will mount this to /var/run/secrets/kubernetes.io/serviceaccount/token
	// but an administrator may have configured it to be mounted elsewhere.
	jwt, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		fmt.Println(fmt.Errorf("unable to read file containing service account token: %w", err))
	}

	params := map[string]interface{}{
		"jwt":  string(jwt),
		"role": cfg.VAULT_ROLE, // the name of the role in Vault that was created with this app's Kubernetes service account bound to it
	}

	// log in to Vault's Kubernetes auth method

	clientAuthMetod := fmt.Sprintf("auth/%s/login", cfg.VAULT_AUTH_METHOD)
	resp, err := client.Logical().Write(clientAuthMetod, params)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to log in with Kubernetes auth: %w", err))
	}
	if resp == nil || resp.Auth == nil || resp.Auth.ClientToken == "" {
		fmt.Println(fmt.Errorf("login response did not return client token"))
	}

	// now you will use the resulting Vault token for making all future calls to Vault
	client.SetToken(resp.Auth.ClientToken)

	// get secrets from Vault
	secret, err := client.Logical().Read(cfg.VAULT_PATH_ENV)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to read secret: %w", err))
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		fmt.Println(fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"]))
	}

	// data map can contain more than one key-value pair, in this case we're just grabbing one of them
	key := cfg.VAULT_ENV
	data[key] = secretEnv

	// put secret to Vault
	_, err = client.Logical().Write(cfg.VAULT_PATH_ENV, data)
	if err != nil {
		fmt.Println(fmt.Errorf("error write secrets to vault: %w", err))
	}

	return err

}

func getToken(cfg config.Config) string {

	type BodyReqStruct struct {
		Client_Id     string `json:"client_id"`
		Client_Secret string `json:"client_secret"`
		Audience      string `json:"audience"`
		Grant_Type    string `json:"grant_type"`
	}

	type BodyResStruct struct {
		Access_Token string `json:"access_token"`
	}

	bodyReqStruct := BodyReqStruct{
		Client_Id:     cfg.API_CLIENT_ID,
		Client_Secret: cfg.API_CLIENT_SECRET,
		Audience:      cfg.API_AUDINCE,
		Grant_Type:    cfg.API_GRANT_TYPE,
	}

	bodyReqJson, err := json.Marshal(bodyReqStruct)
	if err != nil {
		fmt.Println("error:", err)
	}

	req, err := http.NewRequest(cfg.API_METHOD, cfg.API_URL, bytes.NewBuffer(bodyReqJson))
	if err != nil {
		fmt.Printf("Error when create new request:%s", err)
	}

	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Error when sends an HTTP request:%s", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var resBody BodyResStruct
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		fmt.Println("error:", err)
	}

	return resBody.Access_Token
}

func main() {
	Run(DefaultConfig)
}
