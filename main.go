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
	// dynamic metrics
	// - alerts
	thousandAlertDesc = prometheus.NewDesc(
		"thousandeyes_alert",
		"triggered / active alerts for a rule in ThousandEyes.",
		[]string{"test_name", "type", "rule_name", "rule_expression"},
		nil)
	thousandAlertHTMLReachabilitySuccessRatioDesc = prometheus.NewDesc(
		"thousandeyes_alert_html_reachability_ratio",
		"Reachability Success Ratio Gauge defined by: 1 - ViolationCount / MonitorCount ",
		[]string{"test_name", "type", "rule_name", "rule_expression"},
		nil)
	// - bgp tests
	thousandTestBGPReachabilityDesc = prometheus.NewDesc(
		"thousandeyes_test_bgp_reachability_percentage",
		"BGP test ran in ThousandEyes - metric: reachability.",
		[]string{"id", "test_name", "type", "prefix", "country", "monitor_name"},
		nil)
	thousandTestBGPUpdatesDesc = prometheus.NewDesc(
		"thousandeyes_test_bgp_updates",
		"BGP test ran in ThousandEyes - metric: updates.",
		[]string{"id", "test_name", "type", "prefix", "country", "monitor_name"},
		nil)
	thousandTestBGPPathChangesDesc = prometheus.NewDesc(
		"thousandeyes_test_bgp_path_changes",
		"BGP test ran in ThousandEyes - metric: pathChanges.",
		[]string{"id", "test_name", "type", "prefix", "country", "monitor_name"},
		nil)

	// - html tests web
	thousandTestHTMLconnectTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_avg_connect_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: connectTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLDNSTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_avg_dns_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: dnsTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLRedirectsDesc = prometheus.NewDesc(
		"thousandeyes_test_html_num_redirects",
		"HTML test ran in ThousandEyes - metric: NumRedirects.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLreceiveTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_receiveTime_milliseconds",
		"HTML test ran in ThousandEyes - metric: receiveTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLresponseCodeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_response_code",
		"HTML test ran in ThousandEyes - metric: responseCode.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLresponseTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_response_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: responseTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLtotalTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_total_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: totalTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLwaitTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_wait_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: waitTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLwireSizeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_wire_size_byte",
		"HTML test ran in ThousandEyes - metric: wireSize.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)

	// - html tests metrics
	thousandTestHTMLLossDesc = prometheus.NewDesc(
		"thousandeyes_test_html_loss_percentage",
		"HTML test ran in ThousandEyes - metric: loss.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLAvgLatencyDesc = prometheus.NewDesc(
		"thousandeyes_test_html_avg_latency_milliseconds",
		"HTML test ran in ThousandEyes - metric: avgLatency.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLMinLatencyDesc = prometheus.NewDesc(
		"thousandeyes_test_html_min_latency_milliseconds",
		"HTML test ran in ThousandEyes - metric: minLatency.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLMaxLatencyDesc = prometheus.NewDesc(
		"thousandeyes_test_html_max_latency_milliseconds",
		"HTML test ran in ThousandEyes - metric: maxLatency.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	thousandTestHTMLJitterDesc = prometheus.NewDesc(
		"thousandeyes_test_html_jitter_milliseconds",
		"HTML test ran in ThousandEyes - metric: jitter.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
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
	thousandRequestParsingFailMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thousandeyes_parsing_fails",
		Help: "The number request parsing failed.",
	})
	thousandRequestScrapingTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "thousandeyes_scraping_seconds",
		Help: "The number of scraping time in seconds.",
	})
	thousandRequestAPILimitReached = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "thousandeyes_api_request_limit_reached",
		Help: "0 no, 1 hit limit. Request not complete. Tests Details skipped first",
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
	ch <- thousandAlertDesc
	ch <- thousandAlertHTMLReachabilitySuccessRatioDesc

	ch <- thousandTestBGPReachabilityDesc
	ch <- thousandTestBGPUpdatesDesc
	ch <- thousandTestBGPPathChangesDesc

	ch <- thousandTestHTMLLossDesc
	ch <- thousandTestHTMLAvgLatencyDesc
	ch <- thousandTestHTMLMinLatencyDesc
	ch <- thousandTestHTMLMaxLatencyDesc
	ch <- thousandTestHTMLJitterDesc

	ch <- thousandTestHTMLconnectTimeDesc
	ch <- thousandTestHTMLDNSTimeDesc
	ch <- thousandTestHTMLRedirectsDesc
	ch <- thousandTestHTMLreceiveTimeDesc
	ch <- thousandTestHTMLresponseCodeDesc
	ch <- thousandTestHTMLresponseTimeDesc
	ch <- thousandTestHTMLtotalTimeDesc
	ch <- thousandTestHTMLwaitTimeDesc
	ch <- thousandTestHTMLwireSizeDesc
}

func collectAlerts(c collector, ch chan<- prometheus.Metric) {

	thousandRequestsTotalMetric.Inc()

	t, err := c.thousandEyes.GetAlerts()
	if err != nil {
		thousandRequestsFailMetric.Inc()
		return
	}

	a := t.Alert
	for i := range a {

		// alert metrics
		ch <- prometheus.MustNewConstMetric(
			thousandAlertDesc,
			prometheus.GaugeValue,
			float64(a[i].Active),
			a[i].TestName,
			a[i].Type,
			a[i].RuleName,
			a[i].RuleExpression,
		)

		// thousandeyes_parsing_fails
		mC := len(a[i].Monitors)
		if mC == 0 {
			log.Println("Alert Monitor Array is empty - skip thousandeyes_parsing_fails")
		} else {
			rr := 1 - (a[i].ViolationCount / mC)

			ch <- prometheus.MustNewConstMetric(
				thousandAlertHTMLReachabilitySuccessRatioDesc,
				prometheus.GaugeValue,
				float64(rr),
				a[i].TestName,
				a[i].Type,
				a[i].RuleName,
				a[i].RuleExpression,
			)
		}

	}
}
func collectTests(c collector, ch chan<- prometheus.Metric) {

	tBGP, tHTMLm, tHTMLw, err := c.thousandEyes.GetTests()
	thousandRequestsTotalMetric.Inc()

	if err != nil {
		thousandRequestsFailMetric.Inc()
		return
	}

	if err != nil {
		thousandRequestAPILimitReached.Set(1)
	} else {
		thousandRequestAPILimitReached.Set(0)
	}

	for e := range tBGP {

		if len(tBGP[e].Net.BgpMetrics) == 0 {
			log.Println("BGP metrics are emptry for Test:", tBGP[e])
			continue
		}

		log.Println("BGP metrics Test:", tBGP[e])

		for i := range tBGP[e].Net.BgpMetrics {

			fmt.Println(tBGP[e].Net.Test.TestName, " | ", tBGP[e].Net.BgpMetrics[i].Prefix, " | ", tBGP[e].Net.BgpMetrics[i].MonitorName)

			// test BGP metrics
			ch <- prometheus.MustNewConstMetric(
				thousandTestBGPReachabilityDesc,
				prometheus.GaugeValue,
				float64(tBGP[e].Net.BgpMetrics[i].Reachability),
				fmt.Sprintf("%d-%d", tBGP[e].Net.Test.TestID, i),
				tBGP[e].Net.Test.TestName,
				tBGP[e].Net.Test.Type,
				tBGP[e].Net.BgpMetrics[i].Prefix,
				tBGP[e].Net.BgpMetrics[i].CountryID,
				tBGP[e].Net.BgpMetrics[i].MonitorName,
			) /*
				ch <- prometheus.MustNewConstMetric(
					thousandTestBGPUpdatesDesc,
					prometheus.GaugeValue,
					float64(tBGP[e].Net.BgpMetrics[i].Updates),
					string(tBGP[e].Net.Test.TestID),
					tBGP[e].Net.Test.TestName,
					tBGP[e].Net.Test.Type,
					tBGP[e].Net.BgpMetrics[i].Prefix,
					tBGP[e].Net.BgpMetrics[i].CountryID,
					tBGP[e].Net.BgpMetrics[i].MonitorName,
				)
				ch <- prometheus.MustNewConstMetric(
					thousandTestBGPPathChangesDesc,
					prometheus.GaugeValue,
					float64(tBGP[e].Net.BgpMetrics[i].PathChanges),
					fmt.Sprintf("%d",tBGP[e].Net.Test.TestID),
					tBGP[e].Net.Test.TestName,
					tBGP[e].Net.Test.Type,
					tBGP[e].Net.BgpMetrics[i].Prefix,
					tBGP[e].Net.BgpMetrics[i].CountryID,
					tBGP[e].Net.BgpMetrics[i].MonitorName,
				)*/
		}
	}

	for e := range tHTMLm {

		if len(tHTMLm[e].Net.HTTPMetrics) == 0 {
			log.Println("HTML metrics are emptry for Test:", tHTMLm[e])
			continue
		}
		for i := range tHTMLm[e].Net.HTTPMetrics {

			// test HTML metrics
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLAvgLatencyDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].AvgLatency),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLMinLatencyDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].MinLatency),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLMaxLatencyDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].MaxLatency),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLLossDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].Loss),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLJitterDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].Jitter),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
		}
	}

	for e := range tHTMLw {
		if len(tHTMLw[e].Web.HTTPServer) == 0 {
			log.Println("HTML metrics are emptry for Test:", tHTMLw[e])
			continue
		}
		for i := range tHTMLw[e].Web.HTTPServer {

			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLconnectTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ConnectTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLDNSTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].DNSTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLRedirectsDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].NumRedirects),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLreceiveTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ReceiveTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLresponseCodeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ResponseCode),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLresponseTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ResponseTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLtotalTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].TotalTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLwaitTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].WaitTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				thousandTestHTMLwireSizeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].WireSize),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
		}
	}

}

func (c collector) Collect(ch chan<- prometheus.Metric) {

	scrapeStart := time.Now()

	defer func() {
		if r := recover(); r != nil {
			thousandRequestParsingFailMetric.Inc()
			log.Println("Thousand Eyes Parsing Error (", r, ").")
		}
	}()

	collectTests(c, ch)
	collectAlerts(c, ch)

	scrapeElapsed := time.Since(scrapeStart)

	thousandRequestScrapingTime.Add(scrapeElapsed.Seconds())

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	prometheus.MustRegister(thousandRequestParsingFailMetric)
	prometheus.MustRegister(thousandRequestsetRospectionPeriodMetric)
	prometheus.MustRegister(thousandRequestScrapingTime)
	prometheus.MustRegister(thousandRequestAPILimitReached)
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
