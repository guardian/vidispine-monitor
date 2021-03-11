all: vidispine-monitor.amd64

vidispine-monitor.amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vidispine-monitor.amd64

test:
	go test ./...

clean:
	rm -f vidispine-monitor.amd64 vidispine-monitor