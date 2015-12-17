# federation server

Go implementation of [Federation](https://www.stellar.org/developers/learn/concepts/federation.html) protocol server.

## Downloading the tool
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
  * `federation` - prepared statement query to fetch federation results, should return two fields `username` and `account_id`
  * `reverse-federation` - prepared statement query to fetch reverse federation results, should return two fields `username` and `account_id`

#### Example `config.toml`
```toml
domain = "acme.com"
port = 8000

[database]
type = "mysql"
url = "root:@/dbname"

[queries]
federation = "SELECT username, account_id FROM Users WHERE username = ?"
reverse-federation = "SELECT username, account_id FROM Users WHERE account_id = ?"
```

## Usage

```
./federation
```

## Building

[gb](http://getgb.io) is used for building and testing.

```
gb build
```

## Running tests

```
gb test
```
