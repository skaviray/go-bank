# Bank Application using golang

## Database architecture

![alt text](simple-bank.png)


## Install packages required for the development

- Database Migration
  
```bash
brew install golang-migrate
brew install sqlc
migrate create --ext sql -dir db/migration -seq init_schema
go install github.com/golang/mock/mockgen@v1.6.0
export PATH=$PATH:$(go env GOPATH)/bin
```

## Database setup

- setup Database
  
```bash
make postgres-setup
make postgres-start
make createdb
```

- Destroying database

```bash
make postgres-destroy
```

## Starting the service

```bash
make start-server
```