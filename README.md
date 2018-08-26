# GoHTTP

### Description
- This application returns a base64 encoded SHA512 hashed password.
- This application  three resources, each with one method defined
  * POST `/hash`
    - takes a urlencoded form parameter called `password`
    - Returns: text field
  * GET `/stats`
    - Returns: json `{ "Total": 0, "Average": 5000000 }`
    - `Total` is the number of time the /hash endpoint has been hit
    - `Average` is the average time in microseconds that the /hash endpoint took to respond
  * GET `/shutdown` 
    - this endpoint shutsdown the server
    - it makes sure no hashing work is inProgress
    - Returns: Connection Refused
- An error message with an appropriate error code is returned if any issues crop up `{"Error": "some errror message"}`

### Organization
- `rest/endpoint.go` has the Application struct and starts the server
- `rest/endpoint_test.go` tests the application code and makes sure it starts the server
- `handlers/handler.go` has all the endpoint logic
- `handlers/handler_test.go` tests the helper methods and uses httptest to test the handlers

### Setup
```
echo $GOPATH # should be set to $HOME/go
mkdir -p $GOPATH/src/github.com/{{github-user}}
cd $GOPATH/src/github.com/{{github-user}}
git clone https://github.com/rdibari84/GoHTTP.git
```

### Build Code
```
cd $GOPATH/src
go install github.com/rdibari84/GoHTTP/handlers
go install github.com/rdibari84/GoHTTP/rest
```

### Run Unit Tests
- note unit tests use httptest to test api
- also tests concurrent connections
```
cd $GOPATH/src
go test github.com/rdibari84/GoHTTP/handlers
go test github.com/rdibari84/GoHTTP/rest
```

### Run Server
```
cd $GOPATH
bin/rest
```
another way to run
```
cd $GOPATH/src/github.com/{{github-user}}/GoHTTP
go run rest/endpoint.go
```

### Manual Passing Test Commands
```
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X GET http://localhost:8080/stats
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X GET http://localhost:8080/stats
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X GET http://localhost:8080/shutdown
```

### Manual Failing Test Commands
```
# invalid methods
curl -X GET http://localhost:8080/hash
curl -X POST http://localhost:8080/stats
curl -X POST http://localhost:8080/shutdown

# empty form
curl -X POST --data "" http://localhost:8080/hash 
# form instead of an url encoded form
curl -X POST --form "password=angryMonkey" http://localhost:8080/hash
```

