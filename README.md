# federation server

## Installation

```
gb build
```

## Config

The `config.yaml` file must be present in a working directory. Config file should contain following values:

* `port` - server listening port
* `database-type` - database type (sqlite, mysql, postgres)
* `database-url` - url to database connection
* `domain` - domain this federation server represent
* `federation-query` - prepared statement query to fetch federation results, should return two fields `username` and `account_id`
* `reverse-federation-query` - prepared statement query to fetch reverse federation results, should return two fields `username` and `account_id`

#### Example `config.yaml`
```yaml
port: 8000
database-type: mysql
database-url: "root:@/dbname"
domain: acme.com
federation-query: "SELECT username, account_id FROM Users WHERE username = ?"
reverse-federation-query: "SELECT username, account_id FROM Users WHERE account_id = ?"
```

## Usage

```
./federation
```
