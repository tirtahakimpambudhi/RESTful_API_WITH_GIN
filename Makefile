DBString = "user=postgres password=postgres dbname=test sslmode=disable"

MIGRATION_BIN = bin/migration
MAIN_BIN = bin/main
MIGRATEFILE = internal/db/migration/*.sql

migration:
	go build -o $(MIGRATION_BIN) internal/db/migration/migration.go

migrationUpTo:
	$(MIGRATION_BIN) --driver "postgres" --connect $(DBString) --command "up-to" --file "" --dir "internal/db/migration" --version=20240106034532

migrationDownTo:
	$(MIGRATION_BIN) --driver "postgres" --connect $(DBString) --command "down-to" --file "" --dir "internal/db/migration" --version=20240106034532

migrationUp:
	$(MIGRATION_BIN) --driver "postgres" --connect $(DBString) --command "up" --file "" --dir "internal/db/migration"

migrationDown:
	$(MIGRATION_BIN) --driver "postgres" --connect $(DBString) --command "down" --file "" --dir "internal/db/migration"

createUser:
	$(MIGRATION_BIN) --driver "postgres" --connect $(DBString) --command "create" --file "users_table.sql" --dir "internal/db/migration"

createTodos:
	$(MIGRATION_BIN) --driver "postgres" --connect $(DBString) --command "create" --file "todos_list_table.sql" --dir "internal/db/migration"

migrationClean:
	rm -f $(MIGRATEFILE)

build:
	go build -o $(MAIN_BIN) cmd/app/main.go

run:
	$(MAIN_BIN)

clean:
	rm -f $(MIGRATION_BIN) $(MAIN_BIN)
