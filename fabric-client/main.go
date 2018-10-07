package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	ORGNAME     = "org1"
	CHANNELID   = "mychannel"
	USERNAME    = "Admin"
	CHAINCODEID = "testcc"
	FCN         = "invoke"
	CONFIG_PATH = "/opt/gopath/src/fabric-performance-test/fabric-client-go/fixtures/config.yaml"
)

var mainSDK *fabsdk.FabricSDK

func GetHandler(w http.ResponseWriter, r *http.Request) {

	clientChannelContext := mainSDK.ChannelContext(CHANNELID, fabsdk.WithUser(USERNAME), fabsdk.WithOrg(ORGNAME))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	args := [][]byte{[]byte("get"), []byte("a")}

	//data, err := query(client, args)
	query(client, args)
	if err != nil {
		fmt.Println("Failed to query")
		http.Error(w, err.Error(), 500)
		return

	}
	//fmt.Println(string(data))
}

func PutHandler(w http.ResponseWriter, r *http.Request) {

	clientChannelContext := mainSDK.ChannelContext(CHANNELID, fabsdk.WithUser(USERNAME), fabsdk.WithOrg(ORGNAME))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	kv := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))

	args := [][]byte{[]byte("put"), kv, kv}
	err = invoke(client, args)
	if err != nil {
		fmt.Println("Failed to invoke")
		http.Error(w, err.Error(), 500)
		return

	}
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world!!")
}

func main() {
	var err error
	mainSDK, err = fabsdk.New(config.FromFile(CONFIG_PATH))
	if err != nil {
		fmt.Println("Failed to create new SDK: %s", err)
	}

	http.HandleFunc("/gettest", GetHandler)
	http.HandleFunc("/puttest", PutHandler)
	http.HandleFunc("/helloworld", HelloWorldHandler)
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}

func query(client *channel.Client, args [][]byte) ([]byte, error) {
	response, err := client.Query(channel.Request{
		ChaincodeID: CHAINCODEID,
		Fcn:         FCN,
		Args:        args,
	}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
		return nil, err
	}
	//fmt.Printf("--------> response: TransactionID=%s,payload=%s,valid=%s", response.TransactionID, response.Payload, response.TxValidationCode)
	return response.Payload, nil
}

func invoke(client *channel.Client, args [][]byte) error {
	_, err := client.Execute(channel.Request{
		ChaincodeID: CHAINCODEID,
		Fcn:         FCN,
		Args:        args,
	}, channel.WithRetry(retry.DefaultChannelOpts))

	return err
}