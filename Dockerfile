from golang:1.20.2 as build

WORKDIR /usr/src/discat

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o discat .

cmd ["/usr/src/discat/discat"]
