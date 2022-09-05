- [Introduction](#introduction)
- [Installation](#installation)
- [Migration Filename Formate](#migration-filename-formate)
- [Create Migrations With CLI](#create-migrations-with-cli)
- [Create Migrations With Go](#create-migrations-with-go)
- [Create Migration With Makefile](#create-migration-with-makefile)

# Introduction

當在對資料庫操作的相關程式進行測試時可能會發生資料庫版本不同步的情況, 因此在測試前會需要對資料庫進行 migration 的操作

以下總結 database migration 的目的:

- 將 database schema migrate 到不同環境上
- 對 database 異動進行版本控制
- 利用 version control tool(git) 將 database 異動與程式碼版本同步

[golang-migrate](https://github.com/golang-migrate/migrate) 這個 package 主要適用於多個不同的 database engine, 包括 `PostgreSQL`, `MySQL/MariaDB`, `MongoDB` 等

# Installation

macx:

```go
brew install golang-migrate
```

go:

```go
go get github.com/golang-migrate/migrate/v4
```

官方提供了 CLI 與程式運行兩種方式, 使用程式運行更有助於與程式碼的 version control 結合

# Migration Filename Formate

首先在路徑 `migration/` 下準備好 migration 執行檔案, 檔案名稱格式如下:

```go
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

version 可為遞增整數:

```go
1_initialize_schema.down.sql
1_initialize_schema.up.sql
2_add_table.down.sql
2_add_table.up.sql
...
```

或是 timestamp:

```go
1500360784_initialize_schema.down.sql
1500360784_initialize_schema.up.sql
1500445949_add_table.down.sql
1500445949_add_table.up.sql
...
```

`up` 表示在進版時執行, down 反之

# Create Migrations With CLI

創建一個 `users` table:

```shell
migrate create -ext sql -dir db/migrations -seq create_users_table
```

會在 `db/migrations` 資料夾下看到兩個檔案:

- 000001_create_users_table.down.sql
- 000001_create_users_table.up.sql

在 `.up.sql` 檔案中新增 `users` table:

```sql
CREATE TABLE IF NOT EXISTS users(
   user_id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL
);
```

並在 `.down.sql` 檔案中刪除 `users` table:

```sql
DROP TABLE IF EXISTS users;
```

通過 `IF EXISTS/IF NOT EXISTS` 來確保 migration 過程冪等

執行 `migrate`:

```shell
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

透過 `psql example -c "\d users"` 確認剛剛創建的 `users` table:

```sql
                                    Table "public.users"
  Column  |          Type          |                        Modifiers                        
----------+------------------------+---------------------------------------------------------
 user_id  | integer                | not null default nextval('users_user_id_seq'::regclass)
 username | character varying(50)  | not null
 password | character varying(50)  | not null
 email    | character varying(300) | not null
Indexes:
    "users_pkey" PRIMARY KEY, btree (user_id)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)
```

再來確認 reverse migration 也運作正常:

```shell
migrate -database ${POSTGRESQL_URL} -path db/migrations down
```

# Create Migrations With Go

下面舉例 migrations 在 Go 中的使用方式:

Up:

```go
import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:postgres@localhost:5432/example?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
```

Down:

```go
import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:postgres@localhost:5432/example?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
}
```

# Create Migration With Makefile

```makefile
POSTGRESQL_URL="postgres://postgres:postgres@localhost:5432/blog_service?sslmode=disable"

db-migrate-up:
 migrate -source file://db/migrations -database ${POSTGRESQL_URL} up

db-migrate-down:
 migrate -source file://db/migrations -database ${POSTGRESQL_URL} down
 ```