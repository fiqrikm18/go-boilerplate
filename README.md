## Go Boilerplate
This is a simple golang boilerplate using GRPC and REST API. 

### Tech Used
In this project technology used is:

- [Gin](https://gin-gonic.com). 
- [Gorm](https://gorm.io/).

### Prerequisite
- PostgresSQL.
- go v1.19.4.
  
### How To
- Generate certificate for token need using `make certs`. 
- Edit configuration on `config/config.yaml` according to your local machine configuration. 
- Compile Protobuf file into .go using `make pb-gen`, the file will be generated at `internal/model/protobuf` folder. 
- Run source code with make dev on development or `go run main.go`.