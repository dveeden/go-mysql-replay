package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"time"
)

type ReplayStatement struct {
	session int
	epoch   float64
	stmt    string
}

func mysqlsession(c <-chan ReplayStatement, session int, last_stmt_epoch float64) {
	db, err := sql.Open("mysql", "msandbox:msandbox@tcp(127.0.0.1:5709)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	for {
		pkt := <-c
		if last_stmt_epoch != 0.0 {
			sleeptime := time.Duration(pkt.epoch-last_stmt_epoch) * time.Second
			fmt.Printf("Sleeptime: %s\n", sleeptime)
			time.Sleep(sleeptime)
		}
		last_stmt_epoch = pkt.epoch
		fmt.Printf("STATEMENT REPLAY (session: %d): %s\n", session, pkt.stmt)
		_, err := db.Exec(pkt.stmt)
		if err != nil {
			panic(err.Error())
		}
	}
}

func main() {
	fileflag := flag.String("f", "./test.dat", "Path to datafile for replay")
	flag.Parse()

	datFile, err := os.Open(*fileflag)
	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(datFile)
	reader.Comma = '\t'

	pktData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	var startepoch float64 = 0.0
	sessions := make(map[int]chan ReplayStatement)
	for _, stmt := range pktData {
		sessionid, err := strconv.Atoi(stmt[0])
		if err != nil {
			fmt.Println(err)
		}
		epoch, err := strconv.ParseFloat(stmt[1], 64)
		if err != nil {
			fmt.Println(err)
		}
		pkt := ReplayStatement{session: sessionid, epoch: epoch, stmt: stmt[2]}
		if startepoch == 0.0 {
			startepoch = pkt.epoch
		}
		if sessions[pkt.session] != nil {
			sessions[pkt.session] <- pkt
		} else {
			sess := make(chan ReplayStatement)
			sessions[pkt.session] = sess
			go mysqlsession(sessions[pkt.session], pkt.session, startepoch)
			sessions[pkt.session] <- pkt
		}
	}
}
