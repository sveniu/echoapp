FROM golang:alpine
RUN mkdir /build
ADD . /build/
RUN cd /build && go build -o /echoapp
ENTRYPOINT ["/echoapp"]
