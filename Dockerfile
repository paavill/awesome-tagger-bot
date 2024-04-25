FROM golang:1.22.2 as builder
WORKDIR /usr/tagger

COPY . .

RUN apt update
RUN apt install git
RUN apt install openssh-client

RUN chmod 0400 ./rsa
RUN mkdir -p ~/.ssh
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN eval $(ssh-agent -s) && ssh-add ./rsa && git clone git@github.com:paavill/tokens.git

RUN GOOS=linux CGO_ENABLED=1 CC=gcc go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./tagger

FROM ubuntu:20.04
WORKDIR /

RUN apt update
RUN apt install -y ca-certificates

COPY --from=builder /usr/tagger/tagger /bin/tagger
COPY --from=builder /usr/tagger/tokens/awesome_tagger_bot_tg /bin/awesome_tagger_bot_tg

RUN export PATH=$PATH:/bin

ENV BOT_TOKEN_FILENAME=/bin/awesome_tagger_bot_tg

CMD ["/bin/tagger"]