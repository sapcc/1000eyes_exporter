# 1000eyes_exporter

Prometheus exporter for alerts from ThousandEyes.
The port 9350 was chosen because someone already [reserved](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) it for a ThousandEyes exporter that was supposed to be coming soon but has not been published up to November 2018.

# Environment & Arguments

Needs an `ENV VAR "THOUSANDEYES_TOKEN"` set to a valid ThousandEyes token to be able to query.

You can set the period of time it queries into the past with `-retrospectionPeriod`, e.g. `-retrospectionPeriod 12h`. Large values do not make much sense, because we do not get data about when they started or ended. Just that they existed.

# Docker

1. make build

2. Run
-  _Normal Run to get actual alerts firing:_
`docker run --rm -p 9350:9350 -e "THOUSANDEYES_TOKEN=  secret_api_bearer_token  " $(IMAGE):$(VERSION)`

- _Run getting alerts from the past - makes only sense for Check/Debug purpose:_
`docker run --rm -p 9350:9350 -e "THOUSANDEYES_TOKEN=  secret_api_bearer_token " $(IMAGE):$(VERSION) retrospectionPeriod=1800h`