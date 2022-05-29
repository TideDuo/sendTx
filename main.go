package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 交易发起方keystore文件地址
// var fromKeyStoreFile = "/home/blochchain/tide/docker/ethereum/node1/data/keystore/UTC--2022-01-05T11-29-10.256237157Z--3b738d398c0e503704ebb20a68b30362f7640c47"
var fromKeyStoreFile string

// keystore文件对应的密码
var password = ""

// 交易接收方地址，whatever
var toAddress = "0x5f475f85a7c521d857a6c5dde14d3b1ce012cba2"

// http服务地址, 例:http://localhost:8545
// var httpUrl = "http://localhost:8001"
var httpUrl string

// var txcount = rand.Intn(100)
// var txcount = 667

var txcount = 100

var wg sync.WaitGroup

// use out of docker
func main() {
	str1 := "/home/blochchain/tide/docker/ethereum/node"
	srt2 := "/data/keystore/"
	var str3 string
	var keystorePath string
	var httpUrli = "http://localhost:80"
	var N int

	temp, err := ioutil.ReadFile("/home/blochchain/tide/docker/static-nodes.json")
	if err != nil {
		fmt.Println("get N err=", err)
		return
	}
	for _, v := range temp {
		if v == 61 {
			N++
		}
	}
	fmt.Println("N =", N)

	N = 1
	for {
		// for i := N; i >= 1; i-- {
		for i := 1; i <= N; i++ {
			stri := fmt.Sprintf("%d", i)
			keystorePath = fmt.Sprintf(str1 + stri + srt2)
			if i < 10 {
				httpUrl = fmt.Sprintf(httpUrli + "0" + stri)
			} else {
				httpUrl = fmt.Sprintf(httpUrli + stri)
			}
			// fmt.Println(keystorePath)
			// fmt.Println(httpUrl)

			files, err := ioutil.ReadDir(keystorePath)
			if err != nil {
				fmt.Println("get fileDir err=", err)
				return
			}

			// if i == 2 {
			// 	toAddress = "0x131e00ed09db6531acb468e4f620a07ab3c3d790"
			// 	httpUrl = "http://localhost:8002"
			// } else if i == 3 {
			// 	toAddress = "0x3b738d398c0e503704ebb20a68b30362f7640c47"
			// 	httpUrl = "http://localhost:8003"
			// } else {
			// 	toAddress = "0xcdebf5739cbe8549094150e76036f4ffb3129205"
			// 	httpUrl = "http://localhost:8001"
			// }
			for _, f := range files {
				// fmt.Println(f.Name())
				str3 = f.Name()
				// fmt.Println(str3)
				fromKeyStoreFile = fmt.Sprint(keystorePath + str3)
				// fmt.Println(fromKeyStoreFile)
				wg.Add(1)
				// single node
				go TestSendTx(fromKeyStoreFile, toAddress, httpUrl)
			}
			// fromKeyStoreFile = fmt.Sprint(keystorePath + str3)
			// wg.Add(1)
			// go TestSendTx(fromKeyStoreFile, toAddress, httpUrl)
			// wg.Wait()
		}
		wg.Wait()
		// time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

// var count1, count2 int = 1, 2
/*
	以太坊交易发送
*/
func TestSendTx(fromKeyStoreFile string, toAddress string, httpUrl string) {
	// count1++
	// fmt.Println(count1)
	// 创建客户端
	client, err := ethclient.Dial(httpUrl)
	if err != nil {
		fmt.Println("dail err=", err)
		return
	}
	// require.NoError(t, err)

	// var txcount = rand.Intn(100)
	// var txcount = 2000

	// 交易发送方
	// 获取私钥方式一，通过keystore文件
	fromKeystore, err := ioutil.ReadFile(fromKeyStoreFile)
	if err != nil {
		fmt.Println("get privateKey err=", err)
		return
	}
	// fmt.Println(string(fromKeystore))
	// require.NoError(t, err)
	fromKey, err := keystore.DecryptKey(fromKeystore, password)
	if err != nil {
		fmt.Println("decryptKey err=", err)
		return
	}
	fromPrivkey := fromKey.PrivateKey
	fromPubkey := fromPrivkey.PublicKey
	fromAddr := crypto.PubkeyToAddress(fromPubkey)
	// fmt.Println(fromAddr)

	// 获取私钥方式二，通过私钥字符串
	//privateKey, err := crypto.HexToECDSA("私钥字符串")

	// 交易接收方
	// toAddr := common.StringToAddress(toAddress)
	toAddr := common.BytesToAddress([]byte(toAddress))

	// 数量
	amount := big.NewInt(1e18)

	// gasLimit
	var gasLimit uint64 = 21000

	// gasPrice
	var gasPrice *big.Int = big.NewInt(1e9)

	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal("get SuggestGasPrice err=", err)
	// }
	// fmt.Println(gasPrice)

	// chainID
	// chainID, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(chainID)

	// for i := 0; i < 1000; i++ {
	// nonce获取
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		fmt.Println("get nonce err=", err)
		return
	}
	// fmt.Println(nonce)

	// 认证信息组装
	auth := bind.NewKeyedTransactor(fromPrivkey)
	// auth,err := bind.NewTransactor(strings.NewReader(mykey),"111")
	// auth, err := bind.NewKeyedTransactorWithChainID(fromPrivkey, big.NewInt(77))
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount // in wei
	//auth.Value = big.NewInt(100000)     // in wei
	auth.GasLimit = gasLimit // in units
	//auth.GasLimit = uint64(0) // in units
	// auth.GasPrice = gasPrice
	auth.From = fromAddr

	for i := 0; i < txcount; i++ {
		// break
		// 交易创建
		tx := types.NewTransaction(nonce+uint64(i), toAddr, amount, gasLimit, gasPrice, []byte{})

		// 交易签名
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			fmt.Println("signature err=", err)
			return
		}
		// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivkey)
		// signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, fromPrivkey)
		// require.NoError(t, err)

		// 交易发送
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			fmt.Println("SendTransaction err=", err)
		}
		// fmt.Printf("tx sent by %s, tx is: %s\n", fromAddr, signedTx.Hash().Hex()) // tx sent: 0x77006fcb3938f648e2cc65bafd27dec30b9bfbe9df41f78498b9c8b7322a249e
		// 等待挖矿完成
		// bind.WaitMined(context.Background(), client, signedTx)
	}
	// }
	wg.Done()
}
