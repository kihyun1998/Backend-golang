package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	const (
		dbDriver      = "postgres"
		dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
		serverAddress = "0.0.0.0:8080"
	)

	// 서버 생성하려면 db에 연결하고 store 생성해야한다.

	// db 연결
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect : ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalln("CANNOT START SERVER")
	}
}
