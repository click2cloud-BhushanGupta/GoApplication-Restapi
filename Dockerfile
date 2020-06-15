FROM golang:1.14.4-alpine3.11
WORKDIR /app
RUN apk update && apk add --no-cache git


RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/go-sql-driver/mysql
ENV DB_DRIVER=mysql DB_NAME=userinfo Container_name=db_mysql
COPY . .
EXPOSE 8000
RUN go build -o main .

CMD ["/app/main"]

