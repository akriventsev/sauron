package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

var (
	EthClient        *ethclient.Client
	Auth             *bind.TransactOpts
	PrivateKey       *ecdsa.PrivateKey
	ContractAddress  common.Address
	OracleAddress    common.Address
	JobID            string
	EthereumInstance *Main
	Checker          Service
)

func main() {
	forever := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	pause, _ := time.ParseDuration(viper.GetString("check.check_interval"))
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				forever <- struct{}{}
				return
			default:
				Checker.Check()
			}

			time.Sleep(pause)
		}
	}(ctx)

	go func(ctx context.Context) {
		pause, _ := time.ParseDuration(viper.GetString("check.contract_check_interval"))
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				forever <- struct{}{}
				return
			default:
				EthUSDNBallance, err := EthereumInstance.EthUSDNBallance(&bind.CallOpts{})
				if err != nil {
					log.Fatal(err)
				}
				StackedUSDNBallance, err := EthereumInstance.StackedUSDNBallance(&bind.CallOpts{})
				if err != nil {
					log.Fatal(err)
				}
				Checker.StackedUSDNBallance = StackedUSDNBallance
				Checker.EthereumUSDNBallance = EthUSDNBallance
			}
			time.Sleep(pause)
		}
	}(ctx)

	e := echo.New()
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	// Routes
	e.GET("/check", checker)
	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("server.port"))))
	cancel()

}

// Handler
func checker(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Diff float32
	}{
		Checker.DifPercent * 100,
	})
}

func init() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	EthClient, err := ethclient.Dial("https://kovan.infura.io/v3/93f7392cc6444501bff4962617d4be39")
	if err != nil {
		log.Fatal(err)
	}

	PrivateKey, err := crypto.HexToECDSA(viper.GetString("wallet.private_key"))
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	publicKey := PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	Auth = bind.NewKeyedTransactor(PrivateKey)
	Auth.Nonce = big.NewInt(int64(nonce))
	Auth.Value = big.NewInt(0)     // in wei
	Auth.GasLimit = uint64(300000) // in units
	Auth.GasPrice = gasPrice

	ContractAddress = common.HexToAddress(viper.GetString("common.contract"))
	OracleAddress = common.HexToAddress("0x1e04E83652bBF8560F94d55c91536ed1Fc5f66a0")
	JobID = "54dc0b4e270b4ab5b98c0956b925a682"

	EthereumInstance, err = NewMain(ContractAddress, EthClient)
	if err != nil {
		log.Fatal(err)
	}
}
