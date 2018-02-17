FROM golang:alpine

RUN go version
ENV PATH $GOROOT/bin:$PATH
ENV GOPATH /root

# add source
WORKDIR /root/src/github.com/sunmyinf/workplacehub
ADD . .

RUN apk add --update git make
RUN go get -u github.com/golang/dep/cmd/dep
ENV PATH $GOPATH/bin:$PATH
RUN make build

EXPOSE 8010
CMD bin/workplacehub
