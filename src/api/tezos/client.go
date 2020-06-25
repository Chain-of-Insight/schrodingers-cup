package tezos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// MUTEZ is mutez on the tezos network
const MUTEZ = 1000000

var (
	regRPCError = regexp.MustCompile(`\s?\{\s?\"(kind|error)\"\s?:\s?\"[^,]+\"\s?,\s?\"(kind|error)"\s?:\s?\"[^}]+\s?\}\s?`)
)

/*
GoTezos contains a client (http.Client), network contents, and the host of the node. Gives access to
RPC related functions.
*/
type GoTezos struct {
	client client
	host   string
}

/*
RPCError represents and RPC error
*/
type RPCError struct {
	Kind string `json:"kind"`
	Err  string `json:"error"`
}

func (r *RPCError) Error() string {
	return fmt.Sprintf("rpc error (%s): %s", r.Kind, r.Err)
}

/*
RPCErrors represents multiple RPCError(s).s
*/
type RPCErrors []RPCError

type rpcOptions struct {
	Key   string
	Value string
}

type client interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

/*
New returns a pointer to a GoTezos and initializes the library with the host's Tezos netowrk constants.


Parameters:

	host:
		A Tezos node.
*/
func NewTzClient(host string) (*GoTezos, error) {
	gt := &GoTezos{
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		host: cleanseHost(host),
	}

	_, err := gt.Head()
	if err != nil {
		return gt, errors.Wrap(err, "could not initialize library with network constants")
	}

	// constants, err := gt.Constants(block.Hash)
	// if err != nil {
	// 	return gt, errors.Wrap(err, "could not initialize library with network constants")
	// }
	// gt.networkConstants = &constants

	return gt, nil
}

/*
SetClient overrides GoTezos's client. *http.Client satisfies the client interface.

Parameters:

	client:
		A pointer to an http.Client.
*/
func (t *GoTezos) SetClient(client *http.Client) {
	t.client = client
}

/*
Head gets all the information about the head block.

Path:
	/chains/<chain_id>/blocks/head (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks
*/
func (t *GoTezos) Head() (map[string]interface{}, error) {
	var dat map[string]interface{}

	resp, err := t.Get("/chains/main/blocks/head")
	if err != nil {
		return dat, errors.Wrapf(err, "could not get head block")
	}

	err = json.Unmarshal(resp, &dat)
	if err != nil {
		return dat, errors.Wrapf(err, "could not get head block")
	}

	return dat, nil
}

/*
Counter access the counter of a contract, if any.

Path:
	../<block_id>/context/contracts/<contract_id>/counter (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-counter

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	pkh:
		The pkh (address) of the contract for the query.
*/
func (t *GoTezos) Counter(blockhash, pkh string) (int, error) {
	resp, err := t.Get(fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/counter", blockhash, pkh))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}
	var strCounter string
	err = json.Unmarshal(resp, &strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to unmarshal counter")
	}

	counter, err := strconv.Atoi(strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}
	return counter, nil
}

func (t *GoTezos) ForgeOperation(blockhash string, op []byte) (string, error) {
	resp, err := t.Post(fmt.Sprintf("/chains/main/blocks/%s/helpers/forge/operations", blockhash), op)
	if err != nil {
		return "", errors.Wrapf(err, "failed to forge operation")
	}
	var operation string
	err = json.Unmarshal(resp, &operation)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal forged operation")
	}
	return operation, nil
}

/*
InjectionOperation injects an operation in node and broadcast it. Returns the ID of the operation.
The `signedOperationContents` should be constructed using a contextual RPCs from the latest block
and signed by the client. By default, the RPC will wait for the operation to be (pre-)validated
before answering. See RPCs under /blocks/prevalidation for more details on the prevalidation context.
If ?async is true, the function returns immediately. Otherwise, the operation will be validated before
the result is returned. An optional ?chain parameter can be used to specify whether to inject on the
test chain or the main chain.

Path:
	/injection/operation (POST)

Link:
	https/tezos.gitlab.io/api/rpc.html#post-injection-operation

Parameters:

	input:
		Modifies the InjectionOperation RPC query by passing optional URL parameters. Operation is required.
*/
func (t *GoTezos) InjectionOperation(operation string) (string, error) {
	v, err := json.Marshal(operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to inject operation")
	}

	var opts []rpcOptions
	opts = append(opts, rpcOptions{
		"chain",
		"main",
	})

	resp, err := t.Post("/injection/operation", v, opts...)
	if err != nil {
		return "", errors.Wrap(err, "failed to inject operation")
	}

	var opstring string
	err = json.Unmarshal(resp, &opstring)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal operation")
	}

	return opstring, nil
}

/*
ContractStorage gets access the data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/storage (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-storage

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	KT1:
		The contract address.
*/
func (t *GoTezos) ContractStorage(blockhash string, KT1 string) ([]byte, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/storage", blockhash, KT1)
	resp, err := t.Get(query)
	if err != nil {
		return resp, errors.Wrap(err, "could not get storage '%s'")
	}
	return resp, nil
}

func (t *GoTezos) Post(path string, body []byte, opts ...rpcOptions) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", t.host, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *GoTezos) Get(path string, opts ...rpcOptions) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", t.host, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *GoTezos) Delete(path string, opts ...rpcOptions) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", t.host, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *GoTezos) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to complete request")
	}

	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return byts, errors.Wrap(err, "could not read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return byts, fmt.Errorf("response returned code %d with body %s", resp.StatusCode, string(byts))
	}

	err = handleRPCError(byts)
	if err != nil {
		return byts, err
	}

	t.client.CloseIdleConnections()

	return byts, nil
}

func constructQueryParams(req *http.Request, opts ...rpcOptions) {
	q := req.URL.Query()
	for _, opt := range opts {
		q.Add(opt.Key, opt.Value)
	}

	req.URL.RawQuery = q.Encode()
}

func handleRPCError(resp []byte) error {
	if regRPCError.Match(resp) {
		rpcErrors := RPCErrors{}
		err := json.Unmarshal(resp, &rpcErrors)
		if err != nil {
			return nil
		}
		return &rpcErrors[0]
	}

	return nil
}

func cleanseHost(host string) string {
	if len(host) == 0 {
		return ""
	}
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("http://%s", host) //default to http
	}
	return host
}
