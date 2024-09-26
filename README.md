# Go pdfjam
A [pdfjam](https://github.com/rrthomas/pdfjam) wrapper in go with a web server mode and an cli mode with save file dialog

## Requirements
* go
* pdfjam

## Running
Use `go run main.go` to start the http server on port 8080 or `go run main.go [FILE_1] [FILE_2] ... [FILE_N]` to save the joined pdf with a dialog box.

## Building
`go build`
