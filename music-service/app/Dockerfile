FROM golang:alpine3.18

WORKDIR /music-service-src

COPY music-service/app/ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app

EXPOSE 8083

CMD [ "/build/app" ]

