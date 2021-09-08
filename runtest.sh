rm -rf *.log

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

rm -rf *.log