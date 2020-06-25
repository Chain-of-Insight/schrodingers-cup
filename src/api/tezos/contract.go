package tezos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func GetContractStorage(k1 string) (map[string]interface{}, error) {
	var data map[string]interface{}

	url := fmt.Sprintf("https://api.better-call.dev/v1/contract/%s/%s/storage", viper.GetString("TZ_NETWORK"), k1)
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
