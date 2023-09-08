# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

ARG app
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -trimpath -o app cmd/${app}/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/app /bin/app

ENTRYPOINT ["/bin/app"]
