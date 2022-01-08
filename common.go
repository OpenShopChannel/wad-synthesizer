package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/wii-tools/wadlib"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// titleForType retrieves the title ID for the given type.
// For example, assuming an action of "sd" and an app ID of "1", it queries title_ids for sd_title.
func titleForType(action string, appId int) uint64 {
	var column string
	switch action {
	case "sd":
		column = "sd_title"
	case "nand":
		column = "nand_title"
	case "forwarder":
		column = "forwarder_title"
	default:
		log.Fatalf("invalid type %s while running for app ID %d\n", action, appId)
	}

	// Determine the title ID for the given action type.
	var titleId string

	// We sprintf in order to set our column, as defined above.
	query := fmt.Sprintf("SELECT %s FROM title_ids WHERE application_id = $1", column)
	row := pool.QueryRow(ctx, query, appId)
	err := row.Scan(&titleId)
	check(err)

	// Parse to a uint64.
	res, err := strconv.ParseUint(titleId, 16, 64)
	check(err)

	return res
}

// writeToTitlePath writes a file for the given name to the title's path.
func writeForTitle(titleId uint64, filename string, contents []byte) {
	// Ensure the title's directory exists.
	dirPath := fmt.Sprintf("%s/%016x", config.TitlePath, titleId)
	if _, err := os.Stat(dirPath); errors.Is(err, fs.ErrNotExist) {
		// It does not, so we create it!
		err = os.Mkdir(dirPath, 0755)
		check(err)
	}

	filePath := fmt.Sprintf("%s/%s", dirPath, filename)

	err := ioutil.WriteFile(filePath, contents, 0755)
	check(err)
}

// readZip reads the zip file specified by the given UUID.
func readZip(uuid string) []byte {
	filePath := fmt.Sprintf("%s/%s.zip", config.ZipPath, uuid)
	contents, err := ioutil.ReadFile(filePath)
	check(err)

	return contents
}

// updateTicket updates the given title ID with ticket contents and a version.
func updateTicket(titleId uint64, ticket []byte, version int) {
	titleStr := fmt.Sprintf("%016x", titleId)

	// Attempt to insert or update, depending on what is available.
	_, err := pool.Exec(ctx, `INSERT INTO tickets (title_id, ticket, version) VALUES ($1, $2, $3)
		ON CONFLICT(title_id)
		DO UPDATE SET ticket = $2, version = $3`, titleStr, ticket, version)
	check(err)
}

// updateVersion increments the version for this app ID by 1.
func updateVersion(appId int) int {
	var version int

	row := pool.QueryRow(ctx, `UPDATE application
		SET version = version + 1
		WHERE id = $1
		RETURNING version`, appId)
	err := row.Scan(&version)
	check(err)

	return version
}

// createFauxWad loads template files into a usable form for WAD usage.
func createFauxWad(titleId uint64, version int) *wadlib.WAD {
	var fauxWad = new(wadlib.WAD)

	// Create a random title key.
	titleKey := make([]byte, len(fauxWad.Ticket.TitleKey))
	_, err := rand.Read(titleKey)
	// If there is no randomness, something is very wrong.
	check(err)

	// Manipulate a TMD for this title.
	err = fauxWad.LoadTMD(templateTMD)
	check(err)
	fauxWad.TMD.TitleID = titleId
	fauxWad.TMD.TitleVersion = uint16(version)

	// Manipulate a ticket for this title.
	err = fauxWad.LoadTicket(templateTicket)
	check(err)
	fauxWad.Ticket.TitleID = titleId
	fauxWad.Ticket.TitleVersion = uint16(version)

	// Use our new title key.
	var newKey [16]byte
	copy(newKey[:], titleKey[0:16])
	fauxWad.Ticket.UpdateTitleKey(newKey)

	return fauxWad
}
