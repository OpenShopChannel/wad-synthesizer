package main

import (
	"fmt"
	"github.com/wii-tools/wadlib"
)

// importWad overlays the passed WAD onto title files.
func importWad(contents []byte) {
	// Parse the passed WAD.
	wad, err := wadlib.LoadWAD(contents)
	check(err)

	titleId := wad.TMD.TitleID

	// Write this TMD and its certificate chain to disk.
	tmd, err := wad.GetTMD()
	check(err)

	writtenTmd := append(tmd, wad.CertificateChain...)
	writeForTitle(titleId, "tmd", writtenTmd)

	// Next, write the encrypted binary contents to disk.
	// They are named by their content IDs.
	// As they are encrypted, we do not add the .app suffix.
	for _, content := range wad.Data {
		filename := fmt.Sprintf("%016x", content.Record.ID)

		writeForTitle(titleId, filename, content.RawData)
	}

	// Ensure this title's ticket ID is 0.
	wad.Ticket.TicketID = 0

	// Lastly, update the ticket for this title in our database.
	ticket, err := wad.GetTicket()
	check(err)

	updateTicket(titleId, ticket, int(wad.Ticket.TitleVersion))
}
