# TrackSys 2

This is the new version of tracking software for the Digital Production Group.
It tracks current, in-process and historic digitization orders.

### Requirements
* Go 1.19.0+
* Node 17.9+
* exiftool
* Imagemagick
* Kakadu

### Database Notes

The TrackSys2 backend uses a MySQL DB to track everything. The schema is managed by
https://github.com/golang-migrate/migrate and the scripts are in ./backend/db/migrations.

Install the migrate binary on your host system. For OSX, the easiest method is brew. Execute:

`brew install golang-migrate`.

Define your MySQL connection params in an environment variable, like this:

`export TS2DB='mysql://user:password@tcp(host:port)/dbname?x-migrations-table=ts2_schema_migrations'`

Run migrations like this:

`migrate -database ${TS2DB} -path backend/db/migrations up`

Example migrate commads to create a migration and run one:

* `migrate create -ext sql -dir backend/db/migrations -seq update_user_auth`
* `migrate -database ${TS2DB} -path backend/db/migrations/ up`


### Search index notes

Search is powered by Manficore https://manual.manticoresearch.com/Quick_start_guide?client=CONFIG. The
golang integration is handled by github.com/manticoresoftware/manticoresearch-go. On a Mac, the install steps:

* `brew install manticoresearch`
* `brew services start|stop manticoresearch`

The config is located here: /opt/homebrew/etc/manticoresearch/manticore.conf. Note that by default
the manticore binaries are not on the path. Add them with this (put it .zshrc or .bashrc):

`export PATH="/opt/homebrew/opt/manticoresearch/bin:$PATH"`

A sample config is found here:
`./backend/db/index/manticore.conf.sample`

Important: in the config file, each table defination specifies a path. This is where the index will be stored.
This path must exist prior to running the indexer.

The production and staging config templates for the index can be found here:
`/terraform-infrastructure/tracksys-manticore/production/ansible/templates/manticore.conf.template`
`/terraform-infrastructure/tracksys-manticore/staging/ansible/templates/manticore.conf.template`

With the config in the manticore path, the index is regenerated and service restarted with this command:

`indexer --all --rotate`

Once the index has been populated, you can search it with curl:

`curl "http://localhost:9308/sql?query=select%20*%20from%20orders%20where%20match('chance')"`

Or MySQL:
`mysql -h0 -P9306`
`select * from orders where match('chance')`