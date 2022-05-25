FROM golang:latest

WORKDIR / 
COPY . .

RUN go mod download 

EXPOSE 3002

CMD ["go","run","client.go"]