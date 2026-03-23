go test -v -coverpkg=./... -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html
#per mostrare la media di tutti i test
go tool cover -func=coverage.out | tail -n 1 
