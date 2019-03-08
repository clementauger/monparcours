.PHONY: build test clean

GO_BIN ?= go

# some rule to startup db asap.
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
	mkdir -p data
	touch data/monparcours.db
	mkdir -p migrations/{sqlite3,postgres,mysql}/
	make buildassets

buildassets:
	ls client/node_modules > /dev/null || (cd client && npm i)

	mkdir -p client/public/assets/images/
	mkdir -p client/public/assets/files/
	cp -r client/node_modules/leaflet/dist/images/* client/public/assets/images/
	cp -r client/node_modules/typeface-roboto/files/* client/public/assets/files/
	cp -r client/node_modules/leaflet-control-geocoder/dist/images/* client/public/assets/images/
	# cp -r node_modules/leaflet-routing-machine/dist/*\.{png,jpg,jpeg,svg} public/assets/ || echo "continue.."

	client/node_modules/.bin/terser --compress --mangle --safari10 \
	--output client/public/assets/master.min.js -- \
	client/node_modules/leaflet/dist/leaflet.js \
	client/node_modules/leaflet-providers/leaflet-providers.js \
	client/node_modules/leaflet-control-geocoder/dist/Control.Geocoder.js \
	client/node_modules/mithril/mithril.js \
	client/node_modules/moment/min/moment-with-locales.min.js \
	client/node_modules/pikaday/pikaday.js \
	client/public/app/control.coordinates.js

	client/node_modules/.bin/uglifycss \
	client/node_modules/typeface-roboto/index.css \
	client/node_modules/normalize.css/normalize.css \
	client/node_modules/milligram/dist/milligram.min.css \
	client/node_modules/pikaday/css/{pikaday,theme}.css \
	client/node_modules/leaflet/dist/leaflet.css \
	client/node_modules/leaflet-control-geocoder/dist/Control.Geocoder.css \
	client/public/app/control.coordinates.css \
	> client/public/assets/master.min.css

	tree -L 2 client/public

clean:
	rm -fr build


migrate:
	go run main.go m $(name)
	tree migrations

migrateupall:
	go run main.go mu -limit 10000

migrateup:
	go run main.go mu

migratedown:
	go run main.go md

build:
	rm -fr build/*
	mkdir -p build/
	cp app.yml build/app.yml
	cp dbconfig.yml build/dbconfig.yml
	packr -v build --ldflags="-s -w" -o build/monparcours main.go
	tar -czvf build/monparcours.tgz -C build app.yml dbconfig.yml monparcours

test:
	GO_ENV=test_sqlite make test_
test_:
	rm -fr build
	make buildlocal
	mkdir build/data
	touch build/data/monparcours.db
	(cd build; ./monparcours hello)
	(cd build; ./monparcours mu -limit 10000)
	$(eval AKEY := $(shell build/monparcours getkey))
	(cd build; ./monparcours -quiet &)
	(cd client/; (AKEY=$(AKEY) node_modules/.bin/mocha) || killall monparcours)
	killall monparcours

# totally specific to my use case. don t bother.
buildlocal:
	rm -fr build/*
	mkdir -p build/
	(cd deploy; \
		vagrant up || echo "ok"; \
		vagrant rsync || echo "ok"; \
		vagrant ssh -c 'cd /home/vagrant/projects/src/github.com/clementauger/monparcours;packr build --ldflags="-s -w" -o build/monparcours main.go';\
		vagrant scp :/home/vagrant/projects/src/github.com/clementauger/monparcours/build/monparcours ../build/monparcours\
	)
	# vagrant rsync
	cp prod.app.yml build/app.yml
	cp prod.dbconfig.yml build/dbconfig.yml
	tar -czvf build/monparcours.tgz -C build app.yml dbconfig.yml monparcours
	ls -alh build/
	# (cd deploy; vagrant halt)
