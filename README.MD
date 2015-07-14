Пример запуска:


mysql: go run *.go -driver="mysql" -dsn="user:password@/dbname"	

mongo: go run *.go -driver="mongo" -dsn="mongodb://localhost/stat"

couchbase: go run *.go -driver="couchbase" -dsn="couchbase://127.0.0.1:8091"

elasticsearch: go run *.go -driver="elastic" -dsn="127.0.0.1:9200"
