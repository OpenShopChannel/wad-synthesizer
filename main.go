package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Config represents the loaded configuration for this application.s
type Config struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	DB   string `json:"db"`

	TitlePath string `json:"titlePath"`
	ZipPath   string `json:"zipPath"`
}

// conn is the single PostgreSQL connection we will utilize.
var conn *pgx.Conn

// ctx is context.Background under a new and improved name.
var ctx = context.Background()

// config is a shared, global config loaded at application start.
var config Config

// templateTMD is our template TMD.
//go:embed templates/tmd
var templateTMD []byte

// templateTicket is our template ticket.
//go:embed templates/tik
var templateTicket []byte

// templateCerts is our template certificate chain.
//go:embed templates/certs
var templateCerts []byte

func main() {
	if len(os.Args) == 1 || len(os.Args) > 3 {
		fmt.Printf("Usage: %s <type> [app id]\n", os.Args[0])
		fmt.Println("For more information, consult the README.")
		os.Exit(1)
	}

	// Ensure valid criteria.
	switch os.Args[1] {
	case "sd":
	case "nand":
	case "forwarder":
		break
	default:
		fmt.Printf("Invalid generation type %s!\n", os.Args[1])
		os.Exit(1)
	}

	// Load our configuration.
	data, err := ioutil.ReadFile("./config.json")
	check(err)
	err = json.Unmarshal(data, &config)
	check(err)

	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.User, config.Pass, config.Host, config.DB)
	connConfig, err := pgx.ParseConfig(dbString)
	check(err)
	conn, err = pgx.ConnectConfig(ctx, connConfig)

	// Ensure we can connect to PostgreSQL.
	defer conn.Close(ctx)

	// Set up variables useful for scanning in to.
	var appId int
	var titleId string
	var zipUuid string

	// Query all titles, or a specific title.
	if len(os.Args) == 3 {
		// We have a specific app ID passed.
		appId, err = strconv.Atoi(os.Args[2])
		check(err)

		// Obtain necessary metadata.
		row := conn.QueryRow(ctx, `SELECT application.title_id, metadata.file_uuid
		FROM
			application, metadata
		WHERE
			application.id = $1 AND
			application.id = metadata.application_id`, appId)

		err = row.Scan(&titleId, &zipUuid)
		check(err)

		handleTitle(appId, titleId, zipUuid)
	} else {
		// We want to generate this type for all applications.
		rows, _ := conn.Query(ctx, `SELECT application.id, application.title_id, metadata.file_uuid
		FROM
			application, metadata
		WHERE
			application.id = metadata.application_id`)

		if rows.Next() {
			err = rows.Scan(&appId, &titleId, &zipUuid)
			check(err)

			// Handle for all titles!
			handleTitle(appId, titleId, zipUuid)
		}
		// Ensure our query was successful.
		check(rows.Err())
	}
}

// handleTitle determines an appropriate action to perform for the app ID.
func handleTitle(appId int, titleId string, zipUuid string) {
	action := os.Args[1]

	log.Println("handling", action, appId, titleId, zipUuid)

	switch action {
	case "sd":
	case "nand":
	case "forwarder":
		break
	default:
		fmt.Printf("Invalid generation type %s!\n", os.Args[1])
		os.Exit(1)
	}
}

// check ensures everything is okay in this world.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
