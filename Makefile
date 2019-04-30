IMAGE:=hub.global.cloud.sap/monsoon/1000eyes-exporter
VERSION:=v0.0.1

build:
	go get
	go build -o bin/1000eyes-exporter

test:
	GOOS=linux go build -v -o _scratch/test-exporter ./_scratch/
	scp _scratch/test-exporter core@network0.cc.eu-de-1.cloud.sap:

docker:
	docker build -t $(IMAGE):$(VERSION) .

push:
	docker push $(IMAGE):$(VERSION)

