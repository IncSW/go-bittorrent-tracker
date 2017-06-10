package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/IncSW/go-bittorrent-tracker/protocol/http"
	storageBolt "github.com/IncSW/go-bittorrent-tracker/storage/bolt"
	"github.com/boltdb/bolt"
)

var (
	databaseFilename = flag.String("db", "btt.bolt", "The file pathname of the database")
	httpAddress      = flag.String("addr", ":7654", "The address and port of the HTTP server")
)

func main() {
	flag.Parse()

	db, err := bolt.Open(*databaseFilename, 0600, nil)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		db.Close()
		os.Exit(0)
	}()

	server := http.New(storageBolt.New(db))
	if err = server.ListenAndServe(*httpAddress); err != nil {
		panic(err)
	}
}
