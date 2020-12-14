# CLI tool for exporting encrypted Postgres dumps to DigitalOcean Spaces

Export cmd:
	pg_dump -d test -h localhost -p 5436 -U postgres -Ft > backup_test.tar.gz

Restore cmd:
	pg_restore -d test -p 5436 -h localhost  backup_test.tar.gz -c -U postgres
