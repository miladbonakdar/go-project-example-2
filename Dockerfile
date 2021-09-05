FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    ENV=PRODUCTION

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

WORKDIR /src/cmd
RUN go build -o main .

WORKDIR /dist

RUN cp /src/cmd/main .
RUN cp /src/*.env .
RUN mkdir assets
RUN cp /src/assets/*.json ./assets/

EXPOSE 3000

CMD ["/dist/main"]