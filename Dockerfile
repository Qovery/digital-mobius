FROM golang:1.15.5-buster as build

ADD . /digital-mobius
WORKDIR /digital-mobius
RUN go get && go build -o /digital-mobius.bin main.go

FROM debian:buster-slim as run

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
COPY --from=build /digital-mobius.bin /usr/bin/digital-mobius
CMD ["digital-mobius", "recycle"]