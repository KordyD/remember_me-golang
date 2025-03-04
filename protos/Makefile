all: gen

gen:
	protoc -I protos/ protos/**/*.proto --go_out=gens/go/ --go_opt=paths=source_relative --go-grpc_out=gens/go/ --go-grpc_opt=paths=source_relative