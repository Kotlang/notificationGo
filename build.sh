rm -Rf generated
mkdir generated
cd notification-model
protoc --go_out=../generated --go_opt=paths=source_relative \
    --go-grpc_out=../generated --go-grpc_opt=paths=source_relative \
    *.proto
cd ..

cd social-model
protoc --go_out=../generated --go_opt=paths=source_relative \
    --go-grpc_out=../generated --go-grpc_opt=paths=source_relative \
    *.proto
cd ..
go build -mod=mod -o build/ .