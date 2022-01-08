# wad-synthesize
"Look! It's moving! It's alive!"

`wad-synthesize` aims to create WADs suitable to be used with [osc-downloader](https://github.com/OpenShopChannel/osc-downloader) for Homebrew title installation.

Included within the `templates/` directory is an example TMD, ticket and certificate chain utilized
for presets.

## Setup
Examine `config.json.example` for possible settings. A brief synopsis is as follows:
 - `user`, `pass`, `host`, `db`: Utilized to configure the PostgreSQL database metadata is queried from.
 - `zipPath`: Set to a folder containing ZIPs to create with.
   - It is assumed that ZIPs within are named by the UUID on their metadata.
 - `titlePath`: Set to a folder where titles are served via CCS/UCS to the Wii Shop and EC.

## Usage
It is possible to invoke this tool in several ways.
Generically, the form `wad-synthesize <type> [app id]` is accepted.

For example, to work with a specific app ID:
 - `wad-synthesize sd 123` to generate a SD hidden title for app ID 123.
 - `wad-synthesize nand 123` to generate a NAND title for app ID 123.
 - `wad-synthesize forwarder 123` to generate a forwarder for app ID 123.
 - `wad-synthesize all 123` generates all possible types for app ID 123.

You can additionally run operations in bulk, useful for when things have changed.
 - `wad-synthesize sd` generates SD hidden titles for all possible app IDs.
 - `wad-synthesize all` generates all possible types for all possible app IDs.

Currently, only SD generation is implemented. Others will exist in the future.

During runtime, this tool will increment the `version` column within the `application` table.
Previously issued tickets for a given type will retain previous versions until generated again. This is to avoid the user having to update for nothing.