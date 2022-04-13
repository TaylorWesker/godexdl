run-godexdl:
	@go run cmd/godexdl/main.go

run-to_pdf:
	@go run cmd/to_pdf/main.go

godexdl: bin/godexdl

bin/godexdl: cmd/godexdl/main.go
	go build -o $@ $^

to_pdf: bin/to_pdf

bin/to_pdf: cmd/to_pdf/main.go
	go build -o $@ $^
