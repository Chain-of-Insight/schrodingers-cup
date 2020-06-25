package tezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/anchorageoss/tezosprotocol/v2"
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

const UPDATE_VARS_PARAMS = `{"entrypoint":"default","value":{"prim":"Right","args":[{"prim":"Left","args":[{"prim":"Pair","args":[{"prim":"Pair","args":[{"prim":"Pair","args":[{"prim":"Pair","args":[{"int":"%[6]d"},{"int":"%[7]d"}]},{"prim":"Pair","args":[{"int":"%[5]d"},{"int":"%[3]d"}]}]},{"prim":"Pair","args":[{"prim":"Pair","args":[{"int":"%[1]d"},{"int":"%[10]d"}]},{"prim":"Pair","args":[{"int":"%[8]d"},{"int":"%[4]d"}]}]}]},{"prim":"Pair","args":[{"int":"%[9]d"},{"int":"%[2]d"}]}]}]}]}}`

func UpdateGameVars(k1 string, turnDuration int, quorumRatio float64, pointsToWin int, playerStartPts int, rulePassPts int, voteAgainstPts int, ruleFailedPenalty int) (string, error) {
	gt, err := NewTzClient(viper.GetString("TZ_RPC_NODE"))
	if err != nil {
		return "", err
	}

	// get HEAD hash
	head, err := gt.Head()
	if err != nil {
		return "", err
	}

	walletAddr := viper.GetString("TZ_WALLET_PKHSH")

	// get current counter of sender
	counter, err := gt.Counter("head", walletAddr)

	// construct operation
	opStr := `{"branch": "%[1]s",
	"contents": [ { "kind": "transaction",
	            "source": "%[2]s", "fee": "9000",
	            "counter": "%[4]d", "gas_limit": "75000", "storage_limit": "0",
	            "amount": "0",
	            "destination": "%[3]s",
	            "parameters": %[5]s
	} ] }`
	opParams := fmt.Sprintf(UPDATE_VARS_PARAMS, playerStartPts, pointsToWin, rulePassPts, voteAgainstPts, ruleFailedPenalty, 7200, 4, turnDuration, int(quorumRatio), 100)
	opStr = fmt.Sprintf(opStr, head["hash"], walletAddr, k1, counter+1, opParams)

	// forge operation
	forgedOp, err := gt.ForgeOperation("head", []byte(opStr))
	if err != nil {
		return "", err
	}

	// sign forged operation
	signedOp, err := signForgedOperation(forgedOp)
	if err != nil {
		return "", err
	}

	// inject operation
	opHash, err := gt.InjectionOperation(signedOp)
	if err != nil {
		return "", err
	}

	return opHash, nil
}

func signForgedOperation(op string) (string, error) {
	// sign forged op SignOperation
	opBytes, _ := hex.DecodeString(op)
	privKey, _ := tezosprotocol.PrivateKeySeed(viper.GetString("TZ_WALLET_EDSK")).PrivateKey()

	sig, err := SignOperation(opBytes, privKey)
	if err != nil {
		return "", err
	}

	sigBinary, _ := sig.MarshalBinary()
	sigHex := hex.EncodeToString(sigBinary)
	signedOp := op + sigHex
	return signedOp, nil
}
