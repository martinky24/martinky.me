FROM golang:1.20.4-bullseye as builder

WORKDIR /go/bin/

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/server server.go

FROM scratch

ENV MARTINKY_ME_PORT=8080
COPY --from=builder /go/bin/server /bin/server
COPY --from=builder /go/bin/static/ static/
COPY --from=builder /go/bin/templates/ templates/

EXPOSE 8080

ENTRYPOINT ["/bin/server"]
