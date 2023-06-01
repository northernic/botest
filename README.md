# botest-1
电报机器人
功能：扫描网址并通知到群

protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative score.proto
protoc --go_out=. --go_opt=paths=source_relative  *.proto
