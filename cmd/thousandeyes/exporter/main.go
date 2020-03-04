package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/1000eyes_exporter/pkg/thousandeyes"
	"log"
	"net/http"
	"os"
)

var bGetBGP = flag.Bool("GetBGP", false, "-GetBGP true[|false] if you want BGP test data collected")
var bGetHTTP = flag.Bool("GetHTTP", false, "-GetHTTP true[|false] if you want HTTP request test data collected")
var bGetHttpMetrics = flag.Bool("bGetHTTP_METRICS", false, "-GetHttpMetrics true[|false] if you want HTTP routing test data collected")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	log.Println(fmt.Sprintf("We export alerts from %v since scraping\n", thousandeyes.RetrospectionPeriod))
	//var tc = new(thousandeyes. ThousandEyes) //tc.token = *bearerToken
	token := os.Getenv("THOUSANDEYES_TOKEN")

	if token == "" {
		log.Fatal("error: THOUSANDEYES_TOKEN must be set in the Environment Values - it's empty.")
	}

	var c = &thousandeyes.Collector{
		Token : token,
		IsCollectBgp : *bGetBGP,
		IsCollectHttp : *bGetHTTP,
		IsCollectHttpMetrics: *bGetHttpMetrics,
	}
	prometheus.Register(c)

	// make Prometheus client aware of our collector
	thousandeyes.ThousandRequestsetRospectionPeriodMetric.Set(thousandeyes.RetrospectionPeriod.Seconds())

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(
			`<html>
			<head><title>ThousandEyes Alert Exporter</title></head>
			<body>
			<h1>ThousandEyes Alert Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			<p><a href="https://github.com/sapcc/1000eyes_exporter">Git Repository</a></p>
			<p><a href="https://www.thousandeyes.com/">thousandeyes home</a></p>
			</body>
			</html>`))
	})

	// this port has been allocated for a ThousandEyes exporter
	// https://github.com/prometheus/prometheus/wiki/Default-port-allocations
	port := ":9350"
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
