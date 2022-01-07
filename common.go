package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/wii-tools/wadlib"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
)

// appIdChars determines the last 3 characters for the given title ID.
// For example, assuming an application ID of 1, it returns 0x414141, or aAA.
func appIdChars(appId int) uint64 {
	// TODO(spotlightishere): Implement
	return 0x414141
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
	titleStr := strconv.FormatUint(titleId, 16)

	// Attempt to insert or update, depending on what is available.
	_, err := pool.Exec(ctx, `INSERT INTO tickets (title_id, ticket, version) VALUES ($1, $2, $3)
		ON CONFLICT(title_id)
		DO UPDATE SET ticket = $2, version = $3`, titleStr, ticket, version)
	check(err)
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
