# federation server


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
  * `type` - database type (sqlite3, mysql, postgres, cassandra)
  * [sql, required] `url` - url to database connection
  * [cassandra, required] `cluster` - array of cluster IPs (with ports if not listening on default)
  * [cassandra, required] `keyspace` - keyspace to use
  * [cassandra, optional] `proto-version` - version of the native protocol (default `2`)
* `queries`
  * `federation` - Implementation dependent query to fetch federation results, should return either 1 or 3 columns. These columns should be labeled `id`,`memo`,`memo_type` (`memo_type` not supported in Cassandra - it will always be `text`). Memo and memo_type are optional - check [Federation](https://www.stellar.org/developers/learn/concepts/federation.html) docs).
  * `reverse-federation` - Implementation dependent query to fetch reverse federation results, should return one column. This column should be labeled `name`.
* `tls` (only when running HTTPS server)
  * `certificate-file` - a file containing a certificate
  * `private-key-file` - a file containing a matching private key

`memo_type` should be one of the following:
* `id` - then `memo` field should contain unsigned 64-bit integer, please note that this value will be converted to integer so the field should be an integer or a string representing an integer,
* `text` - then `memo` field should contain string, up to 28 characters.
* `hash` - then `memo` field should contain string that is 32bytes base64 encoded.

## Example `config.toml`
In this section you can find config examples for the two main ways of setting up a federation server.

### #1: Every user has their own Stellar account

In case every user owns Stellar account you don't need `memo`. You can simply return `account_id` based on username. Your `queries` section could look like this:

```toml
domain = "acme.com"
port = 8000

[database]
type = "mysql"
url = "root:@/dbname"

[queries]
federation = "SELECT account_id as id FROM Users WHERE username = $1"
reverse-federation = "SELECT username as name FROM Users WHERE account_id = $1"
```


### #2: Single Stellar account for all incoming transactions

If you have a single Stellar account for all incoming transactions you need to use `memo` to check which internal account should receive the payment.

Let's say that your Stellar account ID is: `GAHG6B6QWTC3YNJIKJYUFGRMQNQNEGBALDYNZUEAPVCN2SGIKHTQIKPV` and every user has an `id` and `username` in your database. Then your `queries` section could look like this:

```toml
domain = "acme.com"
port = 8000

[database]
type = "mysql"
url = "root:@/dbname"

[queries]
federation = "SELECT username as memo, 'text' as memo_type, 'GD6WU64OEP5C4LRBH6NK3MHYIA2ADN6K6II6EXPNVUR3ERBXT4AN4ACD' as id FROM Users WHERE username = $1"
reverse-federation = "SELECT username as name FROM Users WHERE account_id = $1"
```

## SQLite sample

`federation-sqlite-sample` is a simple SQLite DB file you can use to test federation server quickly. It contains a single `Users` table with following schema and data:

id | name | accountId
--- | --- | ---
1 | bob | GCW667JUHCOP5Y7KY6KGDHNPHFM4CS3FCBQ7QWDUALXTX3PGXLSOEALY
2 | alice | GCVYGVXNRUUOFYB5OKA37UYBF3W7RK7D6JPNV57FZFYAUU5NKJYZMTK2

It should work out of box with following `config.toml` file:
```toml
domain = "stellar.org"
port = 8000

[database]
type = "sqlite3"
url = "./federation-sqlite-sample"

[queries]
federation = "SELECT accountId as id FROM Users WHERE name = ?"
reverse-federation = "SELECT name FROM Users WHERE accountId = ?"
```

Start the server and then request it:
```
curl "http://localhost:8000/federation?type=name&q=alice*stellar.org"
```

## Usage

```
./federation
```

## Building

[gb](http://getgb.io) is used for building and testing.

Given you have a running golang installation, you can build the server with:

```
gb build
```

After successful completion, you should find `bin/federation` is present in the project directory.

## Running tests

```
gb test
```
