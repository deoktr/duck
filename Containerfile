FROM docker.io/library/golang:1.24 AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY *.go .
RUN go build -o /go/bin/duck

FROM gcr.io/distroless/base-debian12:nonroot

USER nonroot

COPY --from=build /go/bin/duck /usr/bin/duck
ENTRYPOINT [ "/usr/bin/duck" ]
CMD [ "-addr=0.0.0.0:8000" ]
