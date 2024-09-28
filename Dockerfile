FROM golang:1.23.1 as builder
WORKDIR /src
COPY ./ ./
RUN go build

FROM ubuntu:24.04

RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt update && apt install -y texlive-extra-utils

COPY --from=builder /src/go-pdfjam ./app/

CMD ["/app/go-pdfjam"]

