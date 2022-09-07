FROM golang:1.18

RUN mkdir /server

ADD . /server

WORKDIR /server

RUN go mod tidy

RUN cd /server && go build ./main.go

EXPOSE 8800

CMD ["/server/main"]