FROM golang:1.19 as build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY cmd/server/main.go ./

RUN go build -o banana

FROM nvidia/cuda:11.8.0-base-ubuntu22.04

COPY --from=build /build/banana /app/banana

CMD [ "/app/banana" ]