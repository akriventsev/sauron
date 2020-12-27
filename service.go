package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

type Service struct {
	StackedUSDNBallance  *big.Int
	EthereumUSDNBallance *big.Int
	Diference            int64
	DifPercent           float32
}

func (s *Service) updateStackedBallance() {
	stacking := viper.GetString("resources.waves.stacking")
	gateway := viper.GetString("resources.waves.gateway")
	key := viper.GetString("resources.waves.ballance_key")
	resp, err := http.Get(fmt.Sprintf("https://nodes.wavesnodes.com/addresses/data/%s/%s_%s", stacking, key, gateway))
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyStr := buf.String()
	val := gjson.Get(bodyStr, "value")
	n := new(big.Int)
	n, ok := n.SetString(val.String(), 10)
	if !ok {
		log.Fatal("Error")
		return
	}
	s.StackedUSDNBallance = n
	defer resp.Body.Close()
}

func (s *Service) updateEthereumBallance() {
	contract := viper.GetString("resources.ethereum.contract")
	apikey := viper.GetString("resources.ethereum.apikey")

	resp, err := http.Get(fmt.Sprintf("https://api.etherscan.io/api?module=stats&action=tokensupply&contractaddress=%s&apikey=%s", contract, apikey))
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyStr := buf.String()

	val := gjson.Get(bodyStr, "result")

	n := new(big.Int)
	n, ok := n.SetString(val.String(), 10)
	if !ok {
		return
	}

	s.EthereumUSDNBallance = n.Div(n, big.NewInt(1000000000000)) // div by 1000000000000 convert to paulis

	defer resp.Body.Close()
}

func (s *Service) Check() {
	s.updateEthereumBallance()
	s.updateStackedBallance()
	d := big.NewInt(0)
	d.Sub(s.StackedUSDNBallance, s.EthereumUSDNBallance).Abs(d)
	s.Diference = d.Int64()

	dp := big.NewFloat(0.0).SetInt(d)
	dp = dp.Quo(dp, big.NewFloat(0.0).SetInt(s.StackedUSDNBallance))
	s.DifPercent, _ = dp.Float32()
}
