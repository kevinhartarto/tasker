FROM golang:1.23.3

WORKDIR /opt/app

COPY . .

RUN go mod download && go mod verify

RUN go build -v -o /opt/app ./cmd/api

RUN chmod -x /opt/app 

RUN mkdir /opt/app/logs

EXPOSE 3030

CMD [ "/opt/app/api" ]