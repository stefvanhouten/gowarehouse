module example/database

go 1.19

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/sirupsen/logrus v1.9.0
)

require (
	example/logger v0.0.0-00010101000000-000000000000 // indirect
	github.com/caarlos0/env/v6 v6.10.1 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace example/logger => ../logger
