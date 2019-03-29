# 1000eyes_exporter

Prometheus exporter for alerts from ThousandEyes.

Needs an ENV VAR "THOUSANDEYES_TOKEN" set to a valid ThousandEyes token to be able to query.
You can set the period of time it queries into the past with `-retrospectionPeriod`, e.g. `-retrospectionPeriod 12h`. Large values do not make much sense, because we do not get data about when they started or ended. Just that they existed.

The port 9350 was chosen because someone already [reserved](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) it for a ThousandEyes exporter that was supposed to be coming soon but has not been published up to November 2018.