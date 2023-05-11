build-docker:
	docker build -t martinky-site:$$(git rev-parse --short HEAD) .

run:
	go run server.go
