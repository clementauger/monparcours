//md
/*
# monparcours

a small web application so anyone can setup his own protests walk.

It allows users to draw their own path onto maps provided by OpenStreetMap,
later they can share it, print it, save it...


## Install

### for development

  go get -u github.com/clementauger/monparcours
  make install
  make run

the first command is the very classical go command
to fetch a package.

the second will fetch few more dependencies from go/npm and prepare your local tree.

the last command starts the server right away in development.


### for admin

until a packaged binary is provided, you will have to deal with raw source code.

  go get -u github.com/clementauger/monparcours
  make install
  make build
  # (cd build; ./monparcours)

This will generate an archive under `build/`.

Don t forget to adjust the configuration files (`app.yml` / `dbconfig.yml`) before deploying.

### other dependencies

a database server. At that time only `sqlite` and `mysql` are supported.

## Usage

### Commands

```sh
  $ go run main.go -h
  monparcours dev 0.0.0

  usage
     monparcours [action] args...
     monparcours args...

  actions

   getkey|g
    show admin key

   hello
    test db connection

   help|-h|--help
    show help

   migrate|m [name]
    create a new named migration.

   migratedown|md
    revert migrationlittles

   migrateredo|mr
    redo last migration

   migrateskip|mk
    skip next migration

   migratestatus|ms
    show migrations status

   migrateup|mu
    apply migrations

   serve|s
    serve http application
```

### Serving the website

The build does embed all required assets into the binary.

There is a very limited number of files to deploy, one binary and two configuration files.

When serving the app, you will have to deal with rate limiters, cors, some caching.

Some are non configurable, by design, others should be adjusted.

Check the configuration options in the `app.yml` and `dbconfig.yml` they both contain a `sample_env`

### Working with the database

The binary includes neccessary commands to manage
database migrations and schema lifecycle for all stages.

It is based on a well tested library `rubenv/sql-migrate`.

Prepare a new migration with `migrate/m ${name}`, that will generate files under `migrations/{dialect}`.

Use `migrateup/mu`, `migratedown/md` to consume the migrations.

With `migratestatus/ms`, check the status of the migration and which are to be applied.

Check the configuration options in the `app.yml` and `dbconfig.yml` they both contain a `sample_env`

### Working with the frontend

The frontend is a sinple page application,
the build is managed via the Makefile and the commands `buildassets` and `buildfront`.

```sh
  make run # to start developing
  make stop # to end developing
```

### monitoring

the server provides a monitoring interface at port 127.0.0.1:5032

it uses https://github.com/gocraft/health

you  can also monitor the stdout apache logs with a command like

```sh
journalctl --no-tail -f -u monparcours | \
hlogtop \
  -cut=56 \
  -group="asset=.+\.(css|js|png|jpg|gif|ico|woff2?\?.+)$" \
  -group="protest_id=^/protests/[0-9]+" \
  -group="protest_author=^/protests/by_author/.+" \
  -group="wp=(wp-).+"
```

see also the `--no-hostname` flag of the recent `journalctl` versions.

### testing

A small test suite is available under client/test.

Run `make test` to execute them.

It is possible to change the underlying `db` driver used during the tests by adjusting the current stage via the environement variable `GO_ENV`, see the Makefile itself.

### File tree layout

```sh
$ tree -d -L 3
.
├── build                         
├── client
│   ├── node_modules
│   ├── app                           <--- all the frontend is here
│   ├── public                        <--- frontend build
│   │   ├── app                       <--- from ../app folder
│   │   └── assets                    <--- from node modules
│   └── test                          <--- all the frontend testing is here
├── data                              <--- used to store the sqlite database
├── migrations                        <--- migrations are here
│   ├── mysql
│   ├── postgres
│   └── sqlite3
└── server                            <--- application server
    ├── app                           <--- http controllers
    ├── config                        <--- DB config / APP config
    ├── env
    └── model                         <--- database access layer
        ├── mysql
        └── pgsql
  main.go                              <--- main entry point
```


*/

//md
/*
# Licenses

This app is released under `WTFPL`, but it consumes tons of dependencies that
are using more serious licenses.

Don t ever remove the OSM license within the frontend, thank you.

Support those projects if you can.

*/

package main

// goreadme autogen 
