# GoHTTP

#### Setup
```
echo $GOPATH # should be set to $HOME/go
mkdir -p $GOPATH/src/github.com/{{github-user}}
cd $GOPATH/src/github.com/{{github-user}}
git clone https://github.com/rdibari84/GoHTTP.git
```

#### Build Code
```
cd $GOPATH/src
go build github.com/rdibari84/GoHTTP/rest
```

#### Run Unit Tests
- note unit tests use httptest to test api
- also tests concurrent connections
```
cd $GOPATH/src
go test github.com/rdibari84/GoHTTP/rest
```

#### Run Server
```
cd $GOPATH
bin/rest
```
another way to run
```
cd $GOPATH/src/github.com/{{github-user}}/GoHTTP
go run rest/endpoint.go
```

#### Manual Test Commands
```
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X GET http://localhost:8080/stats
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X GET http://localhost:8080/stats
curl -X POST --data "password=angryMonkey" http://localhost:8080/hash
curl -X GET http://localhost:8080/shutdown
```


