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
	// Validate argument length
	if len(os.Args) == 1 || len(os.Args) > 4 {
		fmt.Printf("Usage: %s <action> [arguments]\n", os.Args[0])
		fmt.Println("For more information, consult the README.")
		os.Exit(1)
	}

	// Load our initial configuration.
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

	// Determine if we are importing, or generating.
	action := os.Args[1]
	switch action {
	case "generate":
		handleGenerate()
	case "import":
		handleImport()
	default:
		fmt.Println("Invalid action type specified!")
		fmt.Printf("Usage: %s <action> [arguments]\n", os.Args[0])
		fmt.Println("For more information, consult the README.")
	}
}

// handleImport handles determining logic for imports.
func handleImport() {
	if len(os.Args) == 2 {
		fmt.Println("No WAD file path specified!")
		fmt.Printf("Usage: %s import path/to/title.wad\n", os.Args[0])
		fmt.Println("For more information, consult the README.")
		os.Exit(1)
	}

	path := os.Args[2]
	contents, err := ioutil.ReadFile(path)
	check(err)

	importWad(contents)
	fmt.Println("handled import of", path)
}

// handleGenerate handles determining logic for generation.
func handleGenerate() {
	var err error

	// Determine the current generation type.
	if len(os.Args) == 2 {
		fmt.Println("No generation type specified!")
		fmt.Printf("Usage: %s generate [app id]\n", os.Args[0])
		fmt.Println("For more information, consult the README.")
		os.Exit(1)
	}

	generationType := os.Args[2]

	// Set up variables useful for scanning in to.
	var appId int
	var zipUuid string

	// Query all titles, or a specific title.
	// If we have 4 arguments, we have a specific title passd.
	if len(os.Args) == 4 {
		// We have a specific app ID passed.
		appId, err = strconv.Atoi(os.Args[3])
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

		handleTitle(generationType, appId, zipUuid, version)
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
			handleTitle(generationType, appId, zipUuid, version)
		}
		// Ensure our query was successful.
		check(rows.Err())
	}
}

// handleTitle determines an appropriate action to perform for the app ID.
func handleTitle(generationType string, appId int, zipUuid string, version int) {
	switch generationType {
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
		fmt.Printf("Invalid generation type %s!\n", os.Args[2])
		os.Exit(1)
	}

	fmt.Printf("handled type %s for app ID %d at version %d\n", generationType, appId, version)
}

// check ensures everything is okay in this world.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
