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
  * `federation` - prepared statement query to fetch federation results, should return two fields `username` and `account_id` (and optionally two additional fields: `memo_type` and `memo` - check [Federation](https://www.stellar.org/developers/learn/concepts/federation.html) docs)
  * `reverse-federation` - prepared statement query to fetch reverse federation results, should return two fields `username` and `account_id`

`memo_type` should be one of the following:
* `id` - then `memo` field should contain unsigned 64-bit integer, please note that this value will be converted to integer so the field should be an integer or a string representing an integer,
* `text` - then `memo` field should contain string, up to 28 characters.

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

## Examples

In this section you can find two main scenarios of using federation server.

### #1: Every user owns Stellar account

In case every user owns Stellar account you don't need `memo`. You can simply return `account_id` based on username. Your `queries` section could look like this:

```toml
[queries]
federation = "SELECT username, account_id FROM Users WHERE username = ?"
reverse-federation = "SELECT username, account_id FROM Users WHERE account_id = ?"
```

### #2: Single Stellar account for all incoming transactions

If you have a single Stellar account for all incoming transactions you need to use `memo` to check which internal account should receive the payment.

Let's say that your Stellar account ID is: `GAHG6B6QWTC3YNJIKJYUFGRMQNQNEGBALDYNZUEAPVCN2SGIKHTQIKPV` and every user has an `id` and `username` in your database. Then your `queries` section could look like this:

```toml
[queries]
federation = "SELECT username, 'GAHG6B6QWTC3YNJIKJYUFGRMQNQNEGBALDYNZUEAPVCN2SGIKHTQIKPV' as account_id, 'id' as memo_type, id as memo FROM Users WHERE username = ?"
reverse-federation = "SELECT username, account_id FROM Users WHERE account_id = ?"
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
