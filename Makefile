.PHONY: build test clean

GO_BIN ?= go

# docker database
pgsql:
	docker run --rm --name some-postgres \
		-e POSTGRES_PASSWORD=test \
		-e POSTGRES_USER=test \
		-e POSTGRES_DB=monparcours \
		-p 5432:5432 -d postgres
# psql -h 127.0.0.1 -d monparcours -U test -W

mysql:
	docker run --rm --name some-mysql \
	-e MYSQL_ROOT_PASSWORD=whatever \
	-e MYSQL_PASSWORD=test \
	-e MYSQL_USER=test \
	-e MYSQL_DATABASE=monparcours \
	-p 3306:3306 -d mysql
# mysql -P 3306 -h 127.0.0.1 -utest -ptest

install:
	# prepare dev environment
	(cd client && npm i)
	$(GO_BIN) get -u github.com/gobuffalo/packr
	$(GO_BIN) get -u github.com/maxcnunes/gaper/cmd/gaper
	$(GO_BIN) get -v github.com/rubenv/sql-migrate/...
	$(GO_BIN) get -v github.com/tdewolff/minify
	mkdir -p data
	touch data/monparcours.db
	mkdir -p migrations/{sqlite3,postgres,mysql}/
	make buildfront
	make buildassets

# frontend
stop:
		pkill minify || echo "ok"
		pkill gaper || echo "ok"

run:
		make stop

		gaper -w "server/" -w "main.go" --build-args='-ldflags="-s -w"' --bin-name="build/monparcours" --no-restart-on="error" &

		minify -w --output client/public/app/master.min.js -- \
		client/app/leaflet-color-markers/leaflet-color-markers.js \
		client/app/control.coordinates.js \
		client/app/copy.js \
		client/app/helpers.js \
		client/app/components.js \
		client/app/pages.js \
		client/app/main.js &

		minify -w --output client/public/app/master.min.css -- \
		client/app/css/*.css \
		client/app/themes/orange.css \
		client/app/*.css &

		minify -w --html-keep-document-tags --output client/public/index.html -- \
		client/app/index.html &

buildfront:
		mkdir -p client/public/app/
		mkdir -p client/public/app/leaflet-color-markers/img
		cp -r client/app/leaflet-color-markers/img client/public/app/leaflet-color-markers/
		cp -r client/app/font client/public/app/

		minify -v --output client/public/app/master.min.js -- \
		client/app/leaflet-color-markers/leaflet-color-markers.js \
		client/app/control.coordinates.js \
		client/app/copy.js \
		client/app/helpers.js \
		client/app/components.js \
		client/app/pages.js \
		client/app/main.js

		minify -v --output client/public/app/master.min.css -- \
		client/app/css/*.css \
		client/app/themes/orange.css \
		client/app/*.css

		minify -v --html-keep-document-tags --output client/public/index.html -- \
		client/app/index.html

buildassets:
	ls client/node_modules > /dev/null || (cd client && npm i)

	mkdir -p client/public/assets/images/
	mkdir -p client/public/assets/files/
	cp -r client/node_modules/leaflet/dist/images/* client/public/assets/images/
	cp -r client/node_modules/typeface-roboto/files/* client/public/assets/files/
	cp -r client/node_modules/leaflet-control-geocoder/dist/images/* client/public/assets/images/
	# cp -r node_modules/leaflet-routing-machine/dist/*\.{png,jpg,jpeg,svg} public/assets/ || echo "continue.."

	minify -v --output client/public/assets/master.min.js -- \
	client/node_modules/leaflet/dist/leaflet.js \
	client/node_modules/leaflet-providers/leaflet-providers.js \
	client/node_modules/leaflet-control-geocoder/dist/Control.Geocoder.js \
	client/node_modules/mithril/mithril.js \
	client/node_modules/moment/min/moment-with-locales.min.js \
	client/node_modules/pikaday/pikaday.js

	minify -v --output client/public/assets/master.min.css -- \
	client/node_modules/typeface-roboto/index.css \
	client/node_modules/normalize.css/normalize.css \
	client/node_modules/milligram/dist/milligram.min.css \
	client/node_modules/pikaday/css/{pikaday,theme}.css \
	client/node_modules/leaflet/dist/leaflet.css \
	client/node_modules/leaflet-control-geocoder/dist/Control.Geocoder.css

	tree -L 2 client/public

# migrate
migrate:
	go run main.go m $(name)
	tree migrations

migrateupall:
	go run main.go mu -limit 10000

migrateup:
	go run main.go mu

migratedown:
	go run main.go md

# testing
test:
	GO_ENV=test_sqlite make test_
	
test_all:
	GO_ENV=test_sqlite make test_
	make mysql
	GO_ENV=test_mysql make test_
	make pgsql
	GO_ENV=test_pg make test_

test_:
	make build
	mkdir build/data
	touch build/data/monparcours.db
	(cd build; ./monparcours hello)
	(cd build; ./monparcours md -limit 10000)
	(cd build; ./monparcours mu -limit 10000)
	$(eval AKEY := $(shell build/monparcours getkey))
	(cd build; ./monparcours -quiet &)
	(cd client/; (AKEY=$(AKEY) node_modules/.bin/mocha) || killall monparcours)
	killall monparcours

# build
clean:
	rm -fr build

build:
	rm -fr build/*
	mkdir -p build/
	cp app.yml build/app.yml
	cp dbconfig.yml build/dbconfig.yml
	packr build --ldflags="-s -w" -o build/monparcours main.go
	tar -czvf build/monparcours.tgz -C build app.yml dbconfig.yml monparcours

# totally specific to my use case. don t bother.
buildlocal:
	rm -fr build/*
	mkdir -p build/
	make buildfront
	make buildassets
	sh deploy/build-monparcours.sh
	cp prod.app.yml build/app.yml
	cp prod.dbconfig.yml build/dbconfig.yml
	tar -czvf build/monparcours.tgz -C build app.yml dbconfig.yml monparcours
	ls -alh build/
