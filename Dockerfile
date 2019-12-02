FROM golang:1.13

ENV TINI_VERSION v0.18.0

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /app/knaing_broker ./cmd/broker
RUN rm -rf $GOPATH/pkg/mod

ENTRYPOINT ["/tini", "-s", "--"]
CMD ["/app/knaing_broker"]
