package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/iwanjunaid/pokabox/config"
	"github.com/iwanjunaid/pokabox/manager"
	_ "github.com/lib/pq"
)

func main() {
	var (
		host   = "127.0.0.1"
		port   = 5432
		user   = "postgres"
		pass   = "123456"
		dbname = "postgres"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var (
		groupID = uuid.MustParse("1f830f06-fe7c-450e-b21f-0b8569aad756")
	)

	outboxConfig := config.NewDefaultCommonOutboxConfig(groupID)

	var (
		bootstrapServers = "127.0.0.1:9092"
	)

	kafkaConfig := config.NewCommonKafkaConfig(bootstrapServers)

	manager := manager.New(outboxConfig, kafkaConfig, db)

	manager.Start()
	manager.Await()
}
