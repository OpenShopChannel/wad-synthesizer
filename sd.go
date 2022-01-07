package main

import (
	"github.com/wii-tools/wadlib"
)

// sdTitle returns a usable SD-oriented title ID for the given application ID.
func sdTitle(appId int) uint64 {
	// 00010008-53xxxxxx (Sxxx)
	base := uint64(0x0001008_53000000)
	base += appIdChars(appId)

	return base
}

// generateSD generates a SD title for the given application ID.
// This is an unpacked WAD. Refer to the Wii Shop Channel documentation
// for what this fully entails.
// It additionally inserts an updated ticket to the SOAP table for the title ID.
func generateSD(appId int, zipUuid string, version int) {
	titleId := sdTitle(appId)
	fauxWad := createFauxWad(titleId, version)

	// Read our zip file content.
	zipContent := readZip(zipUuid)

	// Insert our data.
	fauxWad.Data = make([]wadlib.WADFile, 1)
	fauxWad.Data[0] = wadlib.WADFile{
		Record: &fauxWad.TMD.Contents[0],
	}

	// Encrypt our content.
	// This additionally updates our TMD.
	err := fauxWad.UpdateContent(0, zipContent)
	check(err)

	// Write our encrypted content to disk.
	// We only have one - 00000000.
	// As it is encrypted, we do not append .app as an extension.
	writeForTitle(titleId, "00000000", fauxWad.Data[0].RawData)

	// Next, write our TMD to disk.
	tmd, err := fauxWad.GetTMD()
	check(err)

	// The TMD on disk expects to have a certificate chain following it.
	writtenTmd := append(tmd, templateCerts...)
	writeForTitle(titleId, "tmd", writtenTmd)

	// After that, obtain our ticket in a byte form.
	ticket, err := fauxWad.GetTicket()
	check(err)

	// Lastly, update the ticket for this title in our database.
	updateTicket(titleId, ticket, version)
}
