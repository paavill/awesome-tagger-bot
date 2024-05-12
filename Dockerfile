FROM golang:1.22.2 as builder
WORKDIR /usr/tagger

COPY . .

RUN apt update
RUN apt install git

RUN GOOS=linux CGO_ENABLED=0 CC=gcc go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./tagger

FROM ubuntu:20.04 as runner
WORKDIR /

RUN apt update
RUN apt install -y ca-certificates
RUN apt install -y firefox

COPY --from=builder /usr/tagger/tagger /bin/tagger

RUN export PATH=$PATH:/bin

CMD ["/bin/tagger"]