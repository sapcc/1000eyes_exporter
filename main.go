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

var (
	// dynamic metrics
	thousandAlertsDesc = prometheus.NewDesc(
		"thousandeyes_alerts",
		"Alert triggered in ThousandEyes.",
		[]string{"alert_id", "permalink", "test_name", "type"},
		nil)

	// fixed metrics
	thousandRequestsTotalMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thousandeyes_requests_total",
		Help: "The number requests done against ThousandEyes API.",
	})
	thousandRequestsFailMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thousandeyes_requests_fails",
		Help: "The number requests failed against ThousandEyes API.",
	})
	thousandRequestsetRospectionPeriodMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "thousandeyes_retrospection_period_seconds",
		Help: "The number of seconds into the past we query ThousandEyes for.",
	})

	// stuff for 1000 eyes API
	retrospectionPeriod = flag.Duration(
		"retrospectionPeriod",
		0,
		"The number of hours into the past we query ThousandEyes for. You should it just use for Debugging! Syntax: 1800h")

	//bearerToken = flag.String("token", "NOT SET", "Bearer Token of 1oooEyes")
)

type thousandEyes struct {
	token string
}

type collector struct {
	thousandEyes *thousandEyes
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- thousandAlertsDesc
}

func (c collector) Collect(ch chan<- prometheus.Metric) {

	t, err := c.thousandEyes.getAlerts()
	thousandRequestsTotalMetric.Inc()

	if err != nil {
		thousandRequestsFailMetric.Inc()
		return
	}
	a := t.Alert
	for i := range a {
		ch <- prometheus.MustNewConstMetric(
			thousandAlertsDesc,
			prometheus.GaugeValue,
			float64(a[i].Active),
			fmt.Sprintf("%d", a[i].AlertID), a[i].Permalink, a[i].TestName, a[i].Type,
		)
	}
}

func main() {

	flag.Parse()
	fmt.Printf("We export alerts from %v since scraping\n", retrospectionPeriod)
	tc := new(thousandEyes)
	//tc.token = *bearerToken
	tc.token = os.Getenv("THOUSANDEYES_TOKEN")

	if tc.token == "" {
		fmt.Fprintf(os.Stderr, "error: THOUSANDEYES_TOKEN must be set in the Environment Values.\n")
		os.Exit(1)
	}

	collector := &collector{thousandEyes: tc}
	// make Prometheus client aware of our collector
	prometheus.Register(collector)
	prometheus.MustRegister(thousandRequestsTotalMetric)
	prometheus.MustRegister(thousandRequestsFailMetric)
	prometheus.MustRegister(thousandRequestsetRospectionPeriodMetric)
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
