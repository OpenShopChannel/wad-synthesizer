# wad-synthesize
"Look! It's moving! It's alive!"

`wad-synthesize` aims to create WADs suitable to be used with [osc-downloader](https://github.com/OpenShopChannel/osc-downloader) for Homebrew title installation.

It additionally permits importing existing WADs for serving via the Wii Shop Channel.

## Setup
Examine `config.json.example` for possible settings. A brief synopsis is as follows:
 - `user`, `pass`, `host`, `db`: Utilized to configure the PostgreSQL database metadata is queried from.
 - `zipPath`: Set to a folder containing ZIPs to create with.
   - It is assumed that ZIPs within are named by the UUID on their metadata.
 - `titlePath`: Set to a folder where titles are served via CCS/UCS to the Wii Shop and EC.

## Usage
It is possible to invoke this tool in several ways.
Generically, the form `wad-synthesize <action> <type> [app id]` is accepted.

There are two actions: `generate`, and `import`.
Generate synthesizes a WAD and tracks it for downloading, where import unpacks an existing wad for downloading.

### Generation
During runtime, this tool will increment the `version` column within the `application` table.
Previously issued tickets for other types will retain at their present version until generated again. This is to avoid the user having to update for nothing.

For example, to generate with a specific app ID:
 - `wad-synthesize generate sd 123` to generate a SD hidden title for app ID 123.
 - `wad-synthesize generate nand 123` to generate a NAND title for app ID 123.
 - `wad-synthesize generate forwarder 123` to generate a forwarder for app ID 123.
 - `wad-synthesize generate all 123` generates all possible types for app ID 123.

You can additionally run operations in bulk, useful for when things have changed.
 - `wad-synthesize generate sd` generates SD hidden titles for all possible app IDs.
 - `wad-synthesize generate all` generates all possible types for all possible app IDs.

Currently, only SD generation is implemented. Others will exist in the future.

### Importing
This tool will import the title entirely as present within the WAD.
It will unpack its TMD and raw data in a form available to serve to EC.
Lastly, it will insert its ticket into the database, ready to serve.

Usage is simple: `wad-synthesize import path/to/title.wad`