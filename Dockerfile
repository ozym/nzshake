FROM golang:1.6
MAINTAINER "Mark Chadwick" <m.chadwick@gns.cri.nz>

ADD . /go/src/nzshake
RUN go install nzshake

ENTRYPOINT ["/go/bin/nzshake"]
