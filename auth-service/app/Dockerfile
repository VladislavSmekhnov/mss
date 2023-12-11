FROM golang:alpine3.18

WORKDIR /auth-service-src

COPY auth-service/app/ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app

EXPOSE 8082

CMD [ "/build/app" ]