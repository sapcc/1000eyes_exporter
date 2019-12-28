package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	log.Println(fmt.Sprintf("We export alerts from %v since scraping\n", retrospectionPeriod))
	tc := new(thousandEyes)
	//tc.token = *bearerToken
	tc.token = os.Getenv("THOUSANDEYES_TOKEN")

	if tc.token == "" {
		log.Fatal("error: THOUSANDEYES_TOKEN must be set in the Environment Values - it's empty.")
	}

	c := &collector{thousandEyes: tc}
	prometheus.Register(c)
	prometheus.MustRegister(thousandRequestsTotalMetric)
	prometheus.MustRegister(thousandRequestsFailMetric)
	prometheus.MustRegister(thousandRequestParsingFailMetric)
	prometheus.MustRegister(thousandRequestsetRospectionPeriodMetric)
	prometheus.MustRegister(thousandRequestScrapingTime)
	prometheus.MustRegister(thousandRequestAPILimitReached)

	// make Prometheus client aware of our collector
	thousandRequestsetRospectionPeriodMetric.Set(retrospectionPeriod.Seconds())

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`<html>
			<head><title>ThousandEyes Alert Exporter </title></head>
			<body>
			<h1>ThousandEyes Alert Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			<p><a href="https://github.com/sapcc/1000eyes_exporter"</a>Git Repository</p>
			</body>
			</html>`))
	})

	// this port has been allocated for a ThousandEyes exporter
	// https://github.com/prometheus/prometheus/wiki/Default-port-allocations
	port := ":9350"
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
