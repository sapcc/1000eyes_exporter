IMAGE:=keppel.eu-de-1.cloud.sap/ccloud/1000eyes-exporter
VERSION_LATEST:=latest
VERSION:=v0.0.7

build:
	echo ${GOPATH}
	echo ${GOBIN}
	go get -v github.com/prometheus/client_golang/prometheus
	go install ./pkg/thousandeyes
	go build -o ${GOBIN}/thousandeyes-exporter ./cmd/thousandeyes/exporter/main.go

test:
	GOOS=linux go build -v -o _scratch/test-exporter ./_scratch/
	scp _scratch/test-exporter core@network0.cc.eu-de-1.cloud.sap:

docker:
	docker build -t $(IMAGE):$(VERSION) .
	docker build -t $(IMAGE):$(VERSION_LATEST) .

docker-push:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION_LATEST)
