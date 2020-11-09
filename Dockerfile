FROM golang:1.15-alpine as BUILDING_STEP

WORKDIR /go/src/github.com/sapcc/1000eyes_exporter
ENV GOPATH=/go
ENV GOBIN=/go/bin

RUN apk add --no-cache make git
ARG VERSION
ADD . ${WORKDIR}
RUN ls -lisa ${WORKDIR}
RUN make build

FROM alpine:3.8
LABEL maintainer="tilo.geissler@sap.com"
LABEL source_repository="https://github.com/sapcc/1000eyes_exporter"

RUN apk add --no-cache curl
COPY --from=BUILDING_STEP /go/bin/thousandeyes-exporter /usr/local/bin/
RUN ls -lisa /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/thousandeyes-exporter"]
