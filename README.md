`guess-number-grpc` is my first grpc example

# 1 prepare grpc 
```
cd $<project_dir>
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/guessnumber.proto
```

# 2 run server 
```
cd $<project_dir>
go run src/server/main.go
```

# 3 run client
```
cd $<project_dir>
go run src/client/main.go
```

have fun ~V~