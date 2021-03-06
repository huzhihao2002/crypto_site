package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/valyala/fasthttp"
)

var ethereumNetworkMultiplier = float64(1000000000)

func predictEthereumParams() {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://etherscan.io/chart/blocktime?output=csv")
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Ethereum block time: ", err)
		return
	}

	csv := strings.Split(string(resp.Body()), "\n")
	avgBlockTime := make([]float64, len(csv)-2)
	dates := make([]int64, len(csv)-2)
	for i := 1; i < len(csv)-1; i++ {
		strDate := strings.Replace(strings.Replace(strings.Split(csv[i], ",")[1], "\"", "", -1), "\r", "", -1)
		dates[i-1], err = strconv.ParseInt(strDate, 10, 64)
		if err != nil {
			log.Println("Err parsing ethereum block time: ", err)
			return
		}

		strTime := strings.Replace(strings.Replace(strings.Split(csv[i], ",")[2], "\"", "", -1), "\r", "", -1)
		avgBlockTime[i-1], err = strconv.ParseFloat(strTime, 64)
		if err != nil {
			log.Println("Err parsing ethereum block time: ", err)
			return
		}
	}

	req.SetRequestURI("https://etherscan.io/chart/hashrate?output=csv")
	resp = fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Ethereum hashrate: ", err)
		return
	}

	csv = strings.Split(string(resp.Body()), "\n")
	hashrates := make([]float64, len(csv)-2)
	for i := 1; i < len(csv)-1; i++ {
		strHashrate := strings.Replace(strings.Replace(strings.Split(csv[i], ",")[2], "\"", "", -1), "\r", "", -1)
		hashrates[i-1], err = strconv.ParseFloat(strHashrate, 64)
		// hashrates[i-1] *= 1000000000
		if err != nil {
			log.Println("Err parsing ethereum hashrate: ", err)
			return
		}
	}

	req.SetRequestURI("https://etherscan.io/chart/etherprice?output=csv")
	resp = fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		log.Println("Err on requesting Ethereum hashrate: ", err)
		return
	}

	csv = strings.Split(string(resp.Body()), "\n")
	prices := make([]float64, len(csv)-2)
	for i := 1; i < len(csv)-1; i++ {
		strPrice := strings.Replace(strings.Replace(strings.Split(csv[i], ",")[2], "\"", "", -1), "\r", "", -1)
		prices[i-1], err = strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Println("Err parsing ethereum price: ", err)
			return
		}

	}

	fileStats, err := os.Create("data/ethereum_stats.csv")
	if err != nil {
		log.Println("Err on creating ethereum stats file: ", err)
		return
	}

	filePrice, err := os.Create("data/ethereum_price.csv")
	if err != nil {
		log.Println("Err on creating ethereum price file: ", err)
		return
	}

	for i := 0; i < len(hashrates); i++ {
		// fmt.Fprintf(file, "%d,%.40f\n", dates[i], prices[i]/(avgBlockTime[i]*hashrates[i]))
		fmt.Fprintf(fileStats, "%d,%.40f\n", dates[i], (avgBlockTime[i] * hashrates[i]))
		fmt.Fprintf(filePrice, "%d,%.40f\n", dates[i], prices[i])
		// fmt.Fprintf(file, "%d,%.40f\n", dates[i], hashrates[i])
	}

	cmd := exec.Command("java", "-jar",
		"arima/target/garch-1.0-SNAPSHOT-jar-with-dependencies.jar")
	err = cmd.Run()
	if err != nil {
		log.Println("Err on executing ethereum java-arima: ", err)
		return
	}

	bytes, err := ioutil.ReadFile("data/stats_predictions.csv")
	if err != nil {
		log.Println("Err reading ethereum coefficients: ", err)
		return
	}

	arr := strings.Split(strings.Trim(string(bytes), "\n"), "\n")
	log.Println(arr, len(arr))
	newCoefficients := make([]float64, len(arr))
	for i := 0; i < len(arr); i++ {
		newCoefficients[i], _ = strconv.ParseFloat(arr[i], 64)
	}

	bytes, err = ioutil.ReadFile("data/price_predictions.csv")
	if err != nil {
		log.Println("Err reading ethereum coefficients: ", err)
		return
	}

	arr = strings.Split(strings.Trim(string(bytes), "\n"), "\n")
	log.Println(arr, len(arr))
	newPrices := make([]float64, len(arr))
	for i := 0; i < len(arr); i++ {
		newPrices[i], _ = strconv.ParseFloat(arr[i], 64)
	}

	ethereumCoefficients = newCoefficients
	ethereumPrices = newPrices
}
