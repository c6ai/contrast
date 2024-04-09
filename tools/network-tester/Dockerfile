FROM golang:1.21
RUN apt update && apt install -y strace gdb
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
COPY ./*.go ./
RUN go build -o tester tester.go
ENV GOMAXPROCS=1
ENV GOMEMLIMIT=100MiB
CMD strace -f ./tester

# b crypto/x509.readUniqueDirectoryEntries
