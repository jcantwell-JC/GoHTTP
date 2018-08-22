# GoHTTP

#### Setup
```
echo $GOPATH # should be set to $HOME/go
mkdir $GOPATH/src/
cd $GOPATH/src
git clone https://github.com/rdibari84/GoHTTP.git
```

#### Make sure its all working
```
cd $GOPATH/src
go install GoHTTP/stringutil
go install GoHTTP/hello
../bin/hello
```

#### run unit test
```
go test GoHTTP/hash
```
