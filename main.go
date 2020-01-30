package main

import (
	"flag"
	"fmt"
	"offchain-oracles/config"
	"offchain-oracles/server"
	"offchain-oracles/signer"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultConfigFileName = "config.json"
	defaultDbPath         = "db"
	defaultHost           = "127.0.0.1:8080"
)

func main() {
	var host, confFileName, oracleAddress, dbPath string
	flag.StringVar(&host, "host", defaultHost, "set host")
	flag.StringVar(&oracleAddress, "oracleAddress", "", "set oracle address")
	flag.StringVar(&confFileName, "config", defaultConfigFileName, "set config path")
	flag.StringVar(&dbPath, "db", defaultDbPath, "set db path")
	flag.Parse()

	cfg, err := config.Load(confFileName)
	if err != nil {
		panic(err)
	}

	go server.StartServer(host, dbPath)
	go signer.StartSigner(cfg, oracleAddress, dbPath)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Started...")
	<-done
	fmt.Println("Stopped...")
}
