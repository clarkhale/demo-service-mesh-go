FROM golang:1.18

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go get main
RUN go build -o /main

EXPOSE 8081

CMD [ "/main" , "-p=8080"]
