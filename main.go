package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/config"
	events "github.com/iwanjunaid/pokabox/event"
	"github.com/iwanjunaid/pokabox/internal/interfaces/event"
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
	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9092",
		"acks":              "all",
	}

	kafkaConfig := config.NewCommonKafkaConfig(kafkaConfigMap)
	eventHandler := func(e event.Event) {
		switch event := e.(type) {
		case events.PickerStarted:
			fmt.Printf("%v\n", event)
		case events.Picked:
			fmt.Printf("%v\n", event)
		case events.Sent:
			fmt.Printf("%v\n", event)
		case events.StatusChanged:
			fmt.Printf("%v\n", event)
		case events.PickerPaused:
			fmt.Printf("%v\n", event)
		case events.ZombiePickerStarted:
			fmt.Printf("%v\n", event)
		case events.ZombiePicked:
			fmt.Printf("%v\n", event)
		case events.ZombieAcquired:
			fmt.Printf("%v\n", event)
		case events.ZombiePickerPaused:
			fmt.Printf("%v\n", event)
		case events.RemoverStarted:
			fmt.Printf("%v\n", event)
		case events.Removed:
			fmt.Printf("%v\n", event)
		case events.RemoverPaused:
			fmt.Printf("%v\n", event)
		case events.ErrorOccured:
			fmt.Printf("%v\n", event)
		}
	}

	manager, err := manager.New(outboxConfig, kafkaConfig, db)

	if err != nil {
		log.Fatal(err)
	}

	manager.SetEventHandler(eventHandler)
	manager.Start()
	manager.Await()
}
