IMAGE:=hub.global.cloud.sap/monsoon/1000eyes-exporter
VERSION_LATEST:=latest
VERSION:=v0.0.4

build:
	go get
	go build -o bin/1000eyes-exporter

test:
	GOOS=linux go build -v -o _scratch/test-exporter ./_scratch/
	scp _scratch/test-exporter core@network0.cc.eu-de-1.cloud.sap:

docker:
	docker build -t $(IMAGE):$(VERSION) .
	docker build -t $(IMAGE):$(VERSION_LATEST) .

push:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION_LATEST)

