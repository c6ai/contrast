FROM golang:1.21
RUN apt update && apt install -y strace
WORKDIR /app
COPY ./tester.go ./
RUN go build -o tester tester.go
ENV GOMAXPROCS=1
ENV GOMEMLIMIT=100MiB
CMD strace -f ./tester
