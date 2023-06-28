package main

import (
	"cess-faucet/config"
	"cess-faucet/internal/chain"
	"cess-faucet/internal/handler"
	"cess-faucet/logger"
)

func main() {
	logger.LoggerInit()
	config.ConfInit()
	chain.Chain_Init()
	handler.Handler_main()
}
