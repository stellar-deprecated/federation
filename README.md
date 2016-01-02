# federation server

**Alpha version.**

Go implementation of [Federation](https://www.stellar.org/developers/learn/concepts/federation.html) protocol server. This federation server is designed to be dropped in to your existing infrastructure. It can be configured to pull the data it needs out of your existing DB.

## Downloading the server
[Prebuilt binaries](https://github.com/stellar/federation/releases) of the federation server are available on the [releases page](https://github.com/stellar/federation/releases).

| Platform       | Binary file name                                                                         |
|----------------|------------------------------------------------------------------------------------------|
| Mac OSX 64 bit | [federation-darwin-amd64](https://github.com/stellar/federation/releases)      |
| Linux 64 bit   | [federation-linux-amd64](https://github.com/stellar/federation/releases)       |
| Windows 64 bit | [federation-windows-amd64.exe](https://github.com/stellar/federation/releases) |

Alternatively, you can [build](#building) the binary yourself.

## Config

The `config.toml` file must be present in a working directory. Config file should contain following values:

* `domain` - domain this federation server represent
* `port` - server listening port
* `database`
  * `type` - database type (sqlite, mysql, postgres)
  * `url` - url to database connection
* `queries`
  * `federation` - Implementation dependent query to fetch federation results, should return either 1 or 3 columns. These columns should be labeled `id`,`memo`,`type`. Memo and type are optional
  * `reverse-federation` - Implementation dependent query to fetch reverse federation results, should return one column. This column should be labeled `name`.

#### Example `config.toml`
```toml
domain = "acme.com"
port = 8000

[database]
type = "mysql"
url = "root:@/dbname"

[queries]
federation = "SELECT username as memo,"text" as type, 'GD6WU64OEP5C4LRBH6NK3MHYIA2ADN6K6II6EXPNVUR3ERBXT4AN4ACD' as id FROM Users WHERE username = ?"
reverse-federation = "SELECT username as name FROM Users WHERE account_id = ?"
# example for setups with no memo: federation = "SELECT account_id as id FROM Users WHERE username = ?"
```

## Usage

```
./federation
```

## Building

[gb](http://getgb.io) is used for building and testing.

Given you have a running golang installation, you can install this with:

```bash
go get -u github.com/constabulary/gb/...
```

From within the project directory, simply run `gb build`.  After successful
completion, you should find `bin/federation` is present in the project directory.

## Running tests

```
gb test
```


