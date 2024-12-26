- Database Migration
  
```bash
brew install golang-migrate
brew install sqlc
migrate create --ext sql -dir db/migration -seq init_schema
```

```bash
docker-compose up -d
docker logs simple-bank-db-1  -f
docker-compose stop
docker-compose rm -f
rm -rf ~/simple-bank/postgres
docker exec -it simple-bank-db-1  createdb --username=root --owner=root simple_bank
docker exec -it simple-bank-db-1  dropdb simple_bank
```

```test

```
