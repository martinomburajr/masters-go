FROM golang:buster

RUN apt update && apt install -y build-essential r-base

WORKDIR go/src/github.com/martinomburajr/masters-go/

COPY . .

ENTRYPOINT ["sleep", "infinity"]

