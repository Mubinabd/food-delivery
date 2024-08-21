CURRENT_DIR=$(shell pwd)

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}
exp:
	export DBURL='postgres://postgres:1234@localhost:5432/product?sslmode=disable'

mig-run:
	migrate create -ext sql -dir migrations -seq create_table

mig-up:
	migrate -database 'postgres://postgres:1234@localhost:5432/product?sslmode=disable' -path migrations up

mig-down:
	migrate -database 'postgres://postgres:1234@localhost:5432/product?sslmode=disable' -path migrations down

gen-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$PATH:$(go env GOPATH)/bin"
	protoc --go_out=. --go-grpc_out=. protos/*.proto
	protoc --go_out=. --go-grpc_out=. protos/*.proto
	
gen-protoAll:
	protoc --go_out=./ \
	--go-grpc_out=./ \
	protos/*.proto

swag-gen:
	~/go/bin/swag init -g ./api/router.go -o docs force 1	

run:
	go run cmd/main.go