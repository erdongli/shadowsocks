FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

ARG app
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -trimpath -o app cmd/${app}/main.go

FROM alpine

COPY --from=builder /app/app /bin/app

ENTRYPOINT ["/bin/app"]
