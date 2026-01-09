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

* `brew install manticoresoftware/tap/manticoresearch manticoresoftware/tap/manticore-extra`
* `brew services start|stop manticoresearch`

After that, it will be exposed as MySQL running on port 9306; connect with: mysql -P9306 (no password).
The config is located here: /opt/homebrew/etc/manticoresearch/manticore.conf. Note that by default
the manticore binaries are not on the path. Link them to the homebrew bin with something like this:

`ln -s /opt/homebrew/Cellar/manticoresearch/9.2.14-25032816-23296c0f8/bin/indexer manticore_indexer`

The config template for the index can be found here:
`/terraform-infrastructure/tracksys-manticore/production/ansible/templates`

With the config in the manticore path, the index is regenerated with this command:

`manticore_indexer --all`