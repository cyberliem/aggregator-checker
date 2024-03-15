package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ApiResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    GetRoutesResponse `json:"data"`
}

type GetRoutesResponse struct {
	RouteSummary  *RouteSummary `json:"routeSummary"`
	RouterAddress string        `json:"routerAddress"`
}

type (
	RouteSummary struct {
		TokenIn                     string `json:"tokenIn"`
		AmountIn                    string `json:"amountIn"`
		AmountInUSD                 string `json:"amountInUsd"`
		TokenInMarketPriceAvailable bool   `json:"tokenInMarketPriceAvailable"`

		TokenOut                     string `json:"tokenOut"`
		AmountOut                    string `json:"amountOut"`
		AmountOutUSD                 string `json:"amountOutUsd"`
		TokenOutMarketPriceAvailable bool   `json:"tokenOutMarketPriceAvailable"`

		Gas      string `json:"gas"`
		GasPrice string `json:"gasPrice"`
		GasUSD   string `json:"gasUsd"`

		ExtraFee ExtraFee `json:"extraFee"`

		Route [][]Swap `json:"route"`

		Extra RouteExtraData `json:"extra"`
	}

	ExtraFee struct {
		FeeAmount   string `json:"feeAmount"`
		ChargeFeeBy string `json:"chargeFeeBy"`
		IsInBps     bool   `json:"isInBps"`
		FeeReceiver string `json:"feeReceiver"`
	}

	Swap struct {
		Pool              string      `json:"pool"`
		TokenIn           string      `json:"tokenIn"`
		TokenOut          string      `json:"tokenOut"`
		LimitReturnAmount string      `json:"limitReturnAmount"`
		SwapAmount        string      `json:"swapAmount"`
		AmountOut         string      `json:"amountOut"`
		Exchange          string      `json:"exchange"`
		PoolLength        int         `json:"poolLength"`
		PoolType          string      `json:"poolType"`
		PoolExtra         interface{} `json:"poolExtra"`
		Extra             interface{} `json:"extra"`
	}

	ChunkInfo struct {
		AmountIn     string `json:"amountIn"`
		AmountOut    string `json:"amountOut"`
		AmountInUSD  string `json:"amountInUsd"`
		AmountOutUSD string `json:"amountOutUsd"`
	}

	RouteExtraData struct {
		ChunksInfo []ChunkInfo `json:"chunksInfo"`
	}
)

type BuildRouteParams struct {
	RouteSummary RouteSummary `json:"routeSummary"`

	// Sender address of sender wallet
	Sender string `json:"sender"`

	// Recipient address of recipient wallet
	Recipient string `json:"recipient"`

	Deadline          int64  `json:"deadline"`
	SlippageTolerance int64  `json:"slippageTolerance"`
	Referral          string `json:"referral"`
	Source            string `json:"source"`

	// enable gas estimation, default is false
	SkipSimulateTx bool `json:"skipSimulateTx"`
}

func main() {
	// Define the endpoint URL
	endpointURL := "https://aggregator-api.stg.kyberengineering.io/ethereum/api/v1/routes?tokenIn=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&tokenOut=0xdAC17F958D2ee523a2206206994597C13D831ec7&amountIn=10733957955877498808&saveGas=false&gasInclude=true"
	secondEndpointURL := "https://aggregator-api.stg.kyberengineering.io/ethereum/api/v1/route/build"

	// Initialize counters
	successCount := 0
	failureCount := 0
	buildrouteSuccess := 0
	buildRouteFailure := 0
	// Create a ticker for every 30 seconds
	ticker := time.NewTicker(10 * time.Second)

	// Continuously call the endpoint
	for range ticker.C {
		// Make the HTTP request
		response, err := http.Get(endpointURL)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			continue
		}

		// Read the response body
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		// Parse the JSON response
		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			continue
		}

		// Check if the response code is successful
		if apiResponse.Code == 0 {
			// Extract RouteSummary from Data
			routeSummary := apiResponse.Data.RouteSummary

			// Now you can use routeSummary for further processing
			fmt.Println("RouteSummary extracted successfully:", routeSummary)

			// Increment the success count
			successCount++

			params := BuildRouteParams{
				RouteSummary:      *apiResponse.Data.RouteSummary,
				Sender:            "0xa6c883E2dde82FbED20e025BD717a6B7F34f5E6E",
				Recipient:         "0xa6c883E2dde82FbED20e025BD717a6B7F34f5E6E",
				SlippageTolerance: 50,
				Source:            "kyberswap",
				SkipSimulateTx:    false,
			}
			payload, pErr := json.Marshal(params)
			if pErr != nil {
				fmt.Println("Error marshalling payload for second API call:", err)
				return
			}

			secondResponse, err := http.Post(secondEndpointURL, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				fmt.Println("Error making second API request:", err)
				buildRouteFailure++
			} else {
				defer secondResponse.Body.Close()

				// Read the second API response body
				secondBody, err := ioutil.ReadAll(secondResponse.Body)
				if err != nil {
					fmt.Println("Error reading second API response body:", err)
					return
				} else {
					fmt.Println("build route successfully")
					str1 := bytes.NewBuffer(secondBody).String()
					buildrouteSuccess++
					fmt.Println(str1)
				}
			}
		} else {
			// Increment the failure count
			failureCount++
		}

		// Output the counts
		fmt.Printf("Success count: %d Failure count: %d buildRouteSuccessCount: %d buildRouteFailure: %d \n", successCount, failureCount, buildrouteSuccess, buildRouteFailure)
	}
}
