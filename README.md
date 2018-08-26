# GoHTTP

#### Setup
```
echo $GOPATH # should be set to $HOME/go
mkdir -p $GOPATH/src/github.com/{{github-user}}
cd $GOPATH/src/github.com/{{github-user}}
git clone https://github.com/rdibari84/GoHTTP.git
```

#### Make install code
```
cd $GOPATH/src
go install github.com/rdibari84/GoHTTP/rest
```

#### run server
```
cd $GOPATH
bin/rest
```
another way to run
```
cd $GOPATH/src/github.com/{{github-user}}/GoHTTP
go run rest/endpoint.go
```

#### run unit tests
- note unit tests use httptest to test api
- also tests concurrent connections
```
cd $GOPATH/src
go test github.com/rdibari84/GoHTTP/rest
```
