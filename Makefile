run-godexdl:
	@go run src/cmd/godexdl/main.go

run-to_pdf:
	@go run src/cmd/to_pdf/main.go

godexdl: bin/godexdl

bin/godexdl: src/cmd/godexdl/main.go
	go build -o $@ $^

to_pdf: bin/to_pdf

bin/to_pdf: src/cmd/to_pdf/main.go
	go build -o $@ $^