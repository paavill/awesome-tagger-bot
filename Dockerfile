FROM golang:1.22.2 as builder
WORKDIR /usr/tagger

COPY . .

RUN apt update

RUN GOOS=linux CGO_ENABLED=0 CC=gcc go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./tagger

FROM ubuntu:20.04 as runner
WORKDIR /

RUN apt update
RUN apt install -y ca-certificates
RUN apt install -y firefox
RUN apt install -y wget
RUN apt install -y unzip
RUN apt install -y python3-pip
RUN pip install selenium
RUN wget https://github.com/mozilla/geckodriver/releases/download/v0.35.0/geckodriver-v0.35.0-linux32.tar.gz
RUN tar -xvzf geckodriver-v0.35.0-linux32.tar.gz

COPY --from=builder /usr/tagger/tagger /bin/tagger
COPY ./get_news.py ./get_news.py

RUN export PATH=$PATH:/bin

CMD ["/bin/tagger"]