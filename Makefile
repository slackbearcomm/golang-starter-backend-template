postgis:
	docker run -d -p 5438:5432 --name postgis -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgis/postgis:15-3.3

dbshell:
	docker exec -it postgis psql -U root -d gogql

createdb:
	docker exec -it postgis createdb --username=root --owner=root gogql

dropdb:
	docker exec -it postgis dropdb --username=root gogql

migrateup:
	./run-dbmigrate-up.sh

migratedown:
	./run-dbmigrate-down.sh

migratedu:
	make migratedown && make migrateup

dbseed:
	./run-seed.sh

run:
	./run-server.sh

build:
	go build -o main cmd/main.go

.PHONY:
	postgis
	dbshell
	createdb
	dropdb
	migrateup
	migratedown
	migratedu
	dbseed
	run
	build