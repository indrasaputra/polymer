env "local" {
  src = "file://db/schema.sql"
  dev = "docker://postgres/17/dev"
  url = "postgres://postgres:postgres@localhost:54322/postgres?sslmode=disable&search_path=wallet"

  migration {
    dir = "file://db/migrations"
  }
}
