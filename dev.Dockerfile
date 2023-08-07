FROM golang:1.20-alpine
 
WORKDIR /usr/app/src

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod download

CMD ["air", "-c", ".air.toml"]

EXPOSE 5001
