FROM golang:1.11-alpine as BUILDING_STEP

WORKDIR /go/src/github.com/sapcc/1000eyes-exporter
RUN apk add --no-cache make git
ARG VERSION
ADD . .
RUN make build

FROM alpine:3.8
LABEL maintainer="tilo.geissler@sap.com" 

RUN apk add --no-cache curl
COPY --from=BUILDING_STEP /go/src/github.com/sapcc/1000eyes-exporter/bin/ /usr/local/bin/
RUN ls -lisa /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/1000eyes-exporter"]