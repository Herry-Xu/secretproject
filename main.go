package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const chainlinkABI = `[{"inputs":[],"name":"latestRoundData","outputs":[{"internalType":"uint80","name":"roundId","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"}]`

type LastestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func main() {
	loadEnv()
	client := connectToNode()

	chainlinkAddress := os.Getenv("CHAINLINK_CONTRACT_ADDRESS")
	getPrice(client, chainlinkAddress)
}

func getPrice(client *ethclient.Client, contractAddress string) {
	parsedABI, err := abi.JSON(strings.NewReader(chainlinkABI))
	if err != nil {
		log.Fatalf("Failed to parse chainlink ABI: %v", err)
	}

	address := common.HexToAddress(contractAddress)
	contract := bind.NewBoundContract(address, parsedABI, client, client, client)

	var result LastestRoundData

	// Call the latestRoundData function
	err = contract.Call(&bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}, &[]interface{}{
		&result,
	}, "latestRoundData")
	if err != nil {
		log.Fatalf("Failed to get latest round data: %v", err)
	}

	fmt.Printf("RoundId: %d\n", result.RoundId)
	fmt.Printf("Latest Price: %d\n", result.Answer)
	fmt.Printf("Started at: %d\n", result.StartedAt)
	fmt.Printf("Updated at: %d\n", result.UpdatedAt)
	fmt.Printf("Answered In Round: %d\n", result.AnsweredInRound)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func connectToNode() *ethclient.Client {
	nodeUrl := os.Getenv("NODE_URL")
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		log.Fatalf("Failed to connect to node: %v", err)
	}
	return client
}
