# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./views ./views
COPY ./static ./static
RUN go mod download

COPY *.go ./

RUN go build -o /menace

EXPOSE 3000

# ENV NODE_ENV prod

CMD [ "/menace" ]
