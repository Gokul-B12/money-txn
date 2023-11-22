# simplebank

Requirements:
Docker, SQLC, Postgres, table plus, migrate lib, 


Error details:
unknown driver "postgres" (forgotten import?) -> https://askgolang.com/how-to-fix-panic-sql-unknown-driver-postgres-forgotten-import/

migrate create -ext sql -dir db/migration -seq add_sessions
