package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	thousandeyes "github.com/sapcc/1000eyes_exporter/pkg/thousandeyes"
	"log"
	"net/http"
	"os"
)

var evThousandeyesBearerToken = "THOUSANDEYES_BEARER_TOKEN"
var evThousandeyesBasicAuthUser = "THOUSANDEYES_BASIC_AUTH_USER"
var evThousandeyesBasicAuthToken = "THOUSANDEYES_BASIC_AUTH_TOKEN"

var bGetBGP = flag.Bool("GetBGP", false, "-GetBGP=true [true|false (default)] if you want BGP test data collected")
var bGetHTTP = flag.Bool("GetHTTP", false, "-GetHTTP=true [true|false (default)] if you want HTTP request test data collected")
var bGetHttpMetrics = flag.Bool("GetHttpMetrics", false, "-GetHttpMetrics=true [true|false (default)] if you want HTTP routing test data collected")
var retrospectionPeriod = flag.Duration( "RetrospectionPeriodInSec", 0, "give a time going back in Seconds, examples: 10h | 1h10m10s")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	thousandeyes.RetrospectionPeriod = retrospectionPeriod
	thousandeyes.ThousandRequestsetRospectionPeriodMetric.Set(thousandeyes.RetrospectionPeriod.Seconds())
	log.Printf("INFO: History Debug AlertScraping %d", thousandeyes.RetrospectionPeriod)

	isBasicAuth:= false
	user  := ""
	token := os.Getenv(evThousandeyesBearerToken)

	//tbd: refreshToken := os.Getenv("THOUSANDEYES_REFRESH_TOKEN")

	if token == "" {

		user  = os.Getenv(evThousandeyesBasicAuthUser)
		token = os.Getenv(evThousandeyesBasicAuthToken)

		if token == "" || user == "" {
			log.Fatalf("error: %s or the combination of %s and %s must be set in the Environment Values - something is empty.", evThousandeyesBearerToken, evThousandeyesBasicAuthUser, evThousandeyesBasicAuthToken)
		}
		token = os.Getenv(evThousandeyesBasicAuthToken)
		if token == "" {
			log.Fatalf("error: %s or the combination of %s and %s must be set in the Environment Values - something is empty.", evThousandeyesBearerToken, evThousandeyesBasicAuthUser, evThousandeyesBasicAuthToken)
		}
		log.Print("INFO: We use Basic Auth Token for Authentication.")
		isBasicAuth = true
	} else {
		log.Print("INFO: We use Bearer Token for Authentication.")
	}

	var c = &thousandeyes.Collector{
		Token : token,
		User: user,
		IsBasicAuth: isBasicAuth,
		IsCollectBgp : *bGetBGP,
		IsCollectHttp : *bGetHTTP,
		IsCollectHttpMetrics: *bGetHttpMetrics,
	}
	prometheus.Register(c)



	// make Prometheus client aware of our collector
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
