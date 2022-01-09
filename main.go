package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
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

// pool is the PostgreSQL connection pool we will utilize.
var pool *pgxpool.Pool

// ctx is context.Background under a new and improved name.
var ctx = context.Background()

// config is a shared, global config loaded at application start.
var config Config

func main() {
	if len(os.Args) == 1 || len(os.Args) > 3 {
		fmt.Printf("Usage: %s <type> [app id]\n", os.Args[0])
		fmt.Println("For more information, consult the README.")
		os.Exit(1)
	}

	// Ensure valid criteria.
	action := os.Args[1]
	switch action {
	case "all":
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
	connConfig, err := pgxpool.ParseConfig(dbString)
	check(err)
	pool, err = pgxpool.ConnectConfig(ctx, connConfig)

	// Ensure we can connect to PostgreSQL.
	defer pool.Close()

	// Set up variables useful for scanning in to.
	var appId int
	var zipUuid string

	// Query all titles, or a specific title.
	if len(os.Args) == 3 {
		// We have a specific app ID passed.
		appId, err = strconv.Atoi(os.Args[2])
		check(err)

		// Obtain necessary metadata.
		row := pool.QueryRow(ctx, `SELECT metadata.file_uuid
		FROM
			metadata
		WHERE
			metadata.application_id = $1`, appId)

		err = row.Scan(&zipUuid)
		check(err)

		// Obtain the version for this title.
		version := updateVersion(appId)

		handleTitle(action, appId, zipUuid, version)
	} else {
		// We want to generate this type for all applications.
		rows, _ := pool.Query(ctx, `SELECT application.id, metadata.file_uuid
		FROM
			application, metadata
		WHERE
			application.id = metadata.application_id`)

		for rows.Next() {
			err = rows.Scan(&appId, &zipUuid)
			check(err)

			// Obtain the version for this title.
			version := updateVersion(appId)

			// Handle for all titles!
			handleTitle(action, appId, zipUuid, version)
		}
		// Ensure our query was successful.
		check(rows.Err())
	}
}

// handleTitle determines an appropriate action to perform for the app ID.
func handleTitle(action string, appId int, zipUuid string, version int) {
	fmt.Printf("handling action %s for app ID %d at version %d\n", action, appId, version)

	switch action {
	case "all":
		handleTitle("sd", appId, zipUuid, version)
		handleTitle("nand", appId, zipUuid, version)
		handleTitle("forwarder", appId, zipUuid, version)
	case "sd":
		generateSD(appId, zipUuid, version)
	case "nand":
		// not implemented
		break
	case "forwarder":
		// not implemented
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
