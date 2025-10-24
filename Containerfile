FROM docker.io/library/golang:1.25.3-alpine AS build

ARG VERSION=dev

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go .
RUN go build -o /duck \
	-buildvcs=false \
	-trimpath \
	-ldflags "-X 'main.Version=${VERSION}'"

FROM scratch

COPY --from=build /duck /duck
CMD ["/duck"]
