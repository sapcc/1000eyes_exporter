package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	thousandAlertsDesc = prometheus.NewDesc(
		"thousandeyes_alerts",
		"Alert triggered in ThousandEyes.",
		[]string{"alert_id", "permalink", "test_name", "type"},
		nil)
	retrospectionPeriodDesc = prometheus.NewDesc(
		"thousandeyes_retrospection_period_seconds",
		"The number of seconds into the past we query ThousandEyes for.",
		nil,
		nil)
	retrospectionPeriod = flag.Duration(
		"retrospectionPeriod",
		24*time.Hour,
		"The number of hours into the past we query ThousandEyes for")
)

type thousandEyes struct {
	token string
}

type collector struct {
	thousandEyes *thousandEyes
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- thousandAlertsDesc
	ch <- retrospectionPeriodDesc
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	t, err := c.thousandEyes.getAlerts()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(thousandAlertsDesc, err)
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

	ch <- prometheus.MustNewConstMetric(
		retrospectionPeriodDesc,
		prometheus.GaugeValue,
		float64(time.Duration.Seconds(*retrospectionPeriod)),
	)
}

func main() {
	flag.Parse()
	fmt.Printf("We export alerts from %v since scraping\n", retrospectionPeriod)
	tc := new(thousandEyes)
	tc.token = os.Getenv("THOUSANDEYES_TOKEN")
	collector := &collector{thousandEyes: tc}
	// make Prometheus client aware of our collector
	prometheus.Register(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`<html>
			<head><title>ThousandEyes Alert Exporter </title></head>
			<body>
			<h1>ThousandEyes Alert Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	// this port has been allocated for a ThousandEyes exporter
	// https://github.com/prometheus/prometheus/wiki/Default-port-allocations
	port := ":9350"
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
