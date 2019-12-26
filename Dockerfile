FROM golang:buster

RUN apt update && apt install -y build-essential r-base

WORKDIR src/github.com/martinomburajr/masters-go/

RUN Rscript -e 'install.packages("ggplot2")'
RUN Rscript -e 'install.packages("readr")'
RUN Rscript -e 'install.packages("knitr")'
RUN Rscript -e 'install.packages("dplyr")'

COPY . .

RUN go get -d

ENTRYPOINT ["sleep", "infinity"]

