# banana 

Backend that drives the mobile APP HungryHour

## requirements

Golang (version 1.9.3)

postgres

This is a golang APP running on port 8080. 

Data store for the APP is postgres running on port 5432

Follow the steps:

1. Install GoLang
2. Install Postgres
3. go get -u github.com/phassans/banana/
4. create a postgres user 'pshashidhara' with password 'banana123'
5. create a database banana
6. run to setup db
```
psql -h localhost -d banana -U pshashidhara -a -f sql/setup.sql
```
7. cd src/github.com/phassans/banana/
8. Build
```
go build .
```
9. Run: 
```
./banana
```
