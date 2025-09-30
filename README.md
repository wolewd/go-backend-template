### TODO:
- [] HTTP Router & Middleware
- [] Auth Module
- [] Migration Tools schema & seeding
- [] Background Jobs e.g for database cleanup and email sending
- [] Documentation using bruno
- [] Unit testing

### Running The Program:
- Run docker container for Database and S3 Storage
    ```shell
    docker compose up -d
    ```

- Login to S3 Storage with browser using S3 Access Key and Secret Key and create buckets
    ```shell
    localhost:9000
    ```

- Run the backend
    ```shell
    go run cmd/api/main.go
    ```

### Migration Structure & Naming:
```shell
database/migrations/
    0001_init_schema.up.sql
    0001_init_schema.down.sql
    0002_add_users.up.sql
    0002_add_users.down.sql
    ...
```

### Migration Commands:
- Run all migration
    ```shell
    go run cmd/migrate/main.go -action=up
    ```

- Run migration step by step (e.g. 2 steps)
    ```shell
    go run cmd/migrate/main.go -action=up -steps=2
    ```

- Rollback to previous migration
    ```shell
    go run cmd/migrate/main.go -action=down -steps=1
    ```

- Rollback all migration
    ```shell
    go run cmd/migrate/main.go -action=down
    ```

- Drop all schema and tracking
    ```shell
    go run cmd/migrate/main.go -action=drop
    ```

- Check migration version
    ```shell
    go run cmd/migrate/main.go -action=version
    ```

- Use custom migration with override flag
    ```shell
    go run cmd/migrate/main.go -action=up -path=migrations-path
    ```