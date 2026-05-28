module test_app_employees

replace mongox => ../mongox/

replace redisx => ../redisx/

go 1.25.3

require redisx v0.0.0-00010101000000-000000000000

require golang.org/x/sys v0.34.0 // indirect

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/kardianos/service v1.2.4
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/redis/go-redis/v9 v9.17.2 // indirect
	github.com/rs/cors v1.11.1
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver/v2 v2.5.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	mongox v0.0.0-00010101000000-000000000000
)
