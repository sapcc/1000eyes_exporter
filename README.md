# 1000eyes_exporter

Prometheus exporter ot export test metrics and alerts from [ThousandEyes](https://www.thousandeyes.com/).
The port 9350 was chosen because someone already [reserved](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) it for a ThousandEyes exporter that was supposed to be coming soon but has not been published up to November 2018.

# Environment & Arguments

## Environment Settings
Mandatory 

- `ENV VAR "THOUSANDEYES_TOKEN"` 

or 

- `ENV VAR "THOUSANDEYES_BASIC_AUTH_USER"` &&
- `ENV VAR "THOUSANDEYES_BASIC_AUTH_TOKEN"`

set to a valid ThousandEyes token to be able to query.

## Arguments


- `-GetBGP=true [true|false (default)]` if you want BGP test data collected
- `-GetHTTP=true [true|false (default)]` if you want HTTP request test data collected (false is default if not set)
- `-GetHttpMetrics=true [true|false (default)]` if you want HTTP routing test data collected (false is default if not set)
- Just for debugging purpose: `-RetrospectionPeriod` You can set the period of time it queries into the past, e.g. `-RetrospectionPeriod 12h`. Large values do not make much sense, because we do not get data about when they started or ended. Just that they existed.

# Docker

1. make build

2. Run
-  _Normal Run to get actual alerts firing:_

    - Bearer: 
        
        `docker run --rm -p 9350:9350 -e "THOUSANDEYES_BASIC_AUTH_TOKEN=<secret_api_bearer_token>" $(IMAGE):$(VERSION)`
    
    - Basic Auth: 
    
        `docker run --rm -p 9350:9350 -e "THOUSANDEYES_BASIC_AUTH_USER=<secret_api_user>" -e "THOUSANDEYES_BASIC_AUTH_TOKEN=<secret_api_basic_auth_token>" $(IMAGE):$(VERSION)`

-  _Run to get actual alerts firing and Test Results:_

    `docker run --rm -p 9350:9350 -e "THOUSANDEYES_TOKEN=<secret_api_bearer_token>" $(IMAGE):$(VERSION) -GetBGP=true -GetHTTP=true`

- _Run getting alerts from the past - makes only sense for Check/Debug purpose:_

    `docker run --rm -p 9350:9350 -e "THOUSANDEYES_TOKEN=  secret_api_bearer_token " $(IMAGE):$(VERSION) -RetrospectionPeriod=12h`