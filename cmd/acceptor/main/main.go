package main

import (
	"flag"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/labstack/echo"
	"log"
)

var (
	logConfig = flag.String("l", "conf/log.json", "logger config")
	dbConfig  = flag.String("d", "conf/mongodb.json", "database config")
	mqConfig  = flag.String("q", "conf/rabbitmq.json", "queue config")
	port      = flag.String("p", ":8080", "port")
)

func main() {
	flag.Parse()

	e := echo.New()

	logger, err := infrastructure.InitLog(*logConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	db, err := infrastructure.InitDatabase(*dbConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := db.Disconnect(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	queue, err := infrastructure.InitQueue(*mqConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := queue.Close(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()

	app := acceptor.NewAcceptor(e, db, queue, logger, *port)
	log.Fatal(app.Start())
}
