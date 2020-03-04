package thousandeyes

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

var (
	// dynamic metrics
	// - alerts
	ThousandAlertDesc = prometheus.NewDesc(
		"thousandeyes_alert",
		"triggered / active alerts for a rule in ThousandEyes.",
		[]string{"test_name", "type", "rule_name", "rule_expression"},
		nil)
	ThousandAlertHTMLReachabilitySuccessRatioDesc = prometheus.NewDesc(
		"thousandeyes_alert_html_reachability_ratio",
		"Reachability Success Ratio Gauge defined by: 1 - ViolationCount / MonitorCount ",
		[]string{"test_name", "type", "rule_name", "rule_expression"},
		nil)
	// - bgp tests
	ThousandTestBGPReachabilityDesc = prometheus.NewDesc(
		"thousandeyes_test_bgp_reachability_percentage",
		"BGP test ran in ThousandEyes - metric: reachability.",
		[]string{"id", "test_name", "type", "prefix", "country", "monitor_name"},
		nil)
	//ThousandTestBGPUpdatesDesc
	ThousandTestBGPUpdatesDesc = prometheus.NewDesc(
		"thousandeyes_test_bgp_updates",
		"BGP test ran in ThousandEyes - metric: updates.",
		[]string{"id", "test_name", "type", "prefix", "country", "monitor_name"},
		nil)
	ThousandTestBGPPathChangesDesc = prometheus.NewDesc(
		"thousandeyes_test_bgp_path_changes",
		"BGP test ran in ThousandEyes - metric: pathChanges.",
		[]string{"id", "test_name", "type", "prefix", "country", "monitor_name"},
		nil)

	// - html tests web
	ThousandTestHTMLconnectTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_avg_connect_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: connectTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLDNSTimeDesc
	ThousandTestHTMLDNSTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_avg_dns_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: dnsTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLRedirectsDesc
	ThousandTestHTMLRedirectsDesc = prometheus.NewDesc(
		"thousandeyes_test_html_num_redirects",
		"HTML test ran in ThousandEyes - metric: NumRedirects.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLreceiveTimeDesc
	ThousandTestHTMLreceiveTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_receiveTime_milliseconds",
		"HTML test ran in ThousandEyes - metric: receiveTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLresponseCodeDesc
	ThousandTestHTMLresponseCodeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_response_code",
		"HTML test ran in ThousandEyes - metric: responseCode.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	ThousandTestHTMLresponseTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_response_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: responseTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLTotalTimeDesc
	ThousandTestHTMLTotalTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_total_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: totalTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	ThousandTestHTMLwaitTimeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_wait_time_milliseconds",
		"HTML test ran in ThousandEyes - metric: waitTime.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	ThousandTestHTMLwireSizeDesc = prometheus.NewDesc(
		"thousandeyes_test_html_wire_size_byte",
		"HTML test ran in ThousandEyes - metric: wireSize.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)

	// - html tests metrics
	ThousandTestHTMLLossDesc = prometheus.NewDesc(
		"thousandeyes_test_html_loss_percentage",
		"HTML test ran in ThousandEyes - metric: loss.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLAvgLatencyDesc
	ThousandTestHTMLAvgLatencyDesc = prometheus.NewDesc(
		"thousandeyes_test_html_avg_latency_milliseconds",
		"HTML test ran in ThousandEyes - metric: avgLatency.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLMinLatencyDesc
	ThousandTestHTMLMinLatencyDesc = prometheus.NewDesc(
		"thousandeyes_test_html_min_latency_milliseconds",
		"HTML test ran in ThousandEyes - metric: minLatency.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	ThousandTestHTMLMaxLatencyDesc = prometheus.NewDesc(
		"thousandeyes_test_html_max_latency_milliseconds",
		"HTML test ran in ThousandEyes - metric: maxLatency.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)
	//ThousandTestHTMLJitterDesc
	ThousandTestHTMLJitterDesc = prometheus.NewDesc(
		"thousandeyes_test_html_jitter_milliseconds",
		"HTML test ran in ThousandEyes - metric: jitter.",
		[]string{"test_name", "type", "prefix", "country", "agent_name"},
		nil)

	// fixed metrics
	ThousandRequestsTotalMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thousandeyes_requests_total",
		Help: "The number requests done against ThousandEyes API.",
	})
	//ThousandRequestsFailMetric
	ThousandRequestsFailMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thousandeyes_requests_fails",
		Help: "The number requests failed against ThousandEyes API.",
	})
	//ThousandRequestParsingFailMetric
	ThousandRequestParsingFailMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thousandeyes_parsing_fails",
		Help: "The number request parsing failed.",
	})
	//ThousandRequestScrapingTime
	ThousandRequestScrapingTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "thousandeyes_scraping_seconds",
		Help: "The number of scraping time in seconds.",
	})
	//ThousandRequestAPILimitReached
	ThousandRequestAPILimitReached = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "thousandeyes_api_request_limit_reached",
		Help: "0 no, 1 hit limit. Request not complete. Tests Details skipped first",
	})
	ThousandRequestsetRospectionPeriodMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "thousandeyes_retrospection_period_seconds",
		Help: "The number of seconds into the past we query ThousandEyes for.",
	})

	// stuff for 1000 eyes API
	RetrospectionPeriod = flag.Duration(
		"retrospectionPeriod",
		0,
		"The number of hours into the past we query ThousandEyes for. You should it just use for Debugging! Syntax: 1800h")

	//bearerToken = flag.String("Token", "NOT SET", "Bearer Token of 1oooEyes")
)

//type ThousandEyes struct {
//	Token string
//}

type Collector struct {
	//thousandEyes *ThousandEyes
	Token string
	RefreshToken string
	IsCollectBgp bool
	IsCollectHttp bool
	IsCollectHttpMetrics bool
}

func (t *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ThousandAlertDesc
	ch <- ThousandAlertHTMLReachabilitySuccessRatioDesc

	ch <- ThousandTestBGPReachabilityDesc
	ch <- ThousandTestBGPUpdatesDesc
	ch <- ThousandTestBGPPathChangesDesc

	ch <- ThousandTestHTMLLossDesc
	ch <- ThousandTestHTMLAvgLatencyDesc
	ch <- ThousandTestHTMLMinLatencyDesc
	ch <- ThousandTestHTMLMaxLatencyDesc
	ch <- ThousandTestHTMLJitterDesc

	ch <- ThousandTestHTMLconnectTimeDesc
	ch <- ThousandTestHTMLDNSTimeDesc
	ch <- ThousandTestHTMLRedirectsDesc
	ch <- ThousandTestHTMLreceiveTimeDesc
	ch <- ThousandTestHTMLresponseCodeDesc
	ch <- ThousandTestHTMLresponseTimeDesc
	ch <- ThousandTestHTMLTotalTimeDesc
	ch <- ThousandTestHTMLwaitTimeDesc
	ch <- ThousandTestHTMLwireSizeDesc

}
func addStaticMetrics(ch chan<- prometheus.Metric){
	ch <- ThousandRequestsTotalMetric
	ch <- ThousandRequestsFailMetric
	ch <- ThousandRequestParsingFailMetric
	ch <- ThousandRequestsetRospectionPeriodMetric
	ch <- ThousandRequestScrapingTime
	ch <- ThousandRequestAPILimitReached
}

func collectAlerts(c Collector, ch chan<- prometheus.Metric) {

	

	ThousandRequestsTotalMetric.Inc()

	t, bHitRateLimit, bError  :=  c.GetAlerts()

	// hint for the limit
	if bHitRateLimit {
		ThousandRequestAPILimitReached.Set(1)
	} else {
		ThousandRequestAPILimitReached.Set(0)
	}

	// if we hit the api request limit in between we skip the test details
	if bError {
		ThousandRequestsFailMetric.Inc()
		return
	}

	a := t.Alert
	for i := range a {

		// alert metrics
		ch <- prometheus.MustNewConstMetric(
			ThousandAlertDesc,
			prometheus.GaugeValue,
			float64(a[i].Active),
			a[i].TestName,
			a[i].Type,
			a[i].RuleName,
			a[i].RuleExpression,
		)

		// skip thousandeyes_parsing_fails for non BGP alerts
		mC := len(a[i].Monitors)
		if mC == 0 {
			log.Println("INFO: Alert Monitor Array is empty (only for BGP)")
		} else {
			rr := 1 - (a[i].ViolationCount / mC)

			ch <- prometheus.MustNewConstMetric(
				ThousandAlertHTMLReachabilitySuccessRatioDesc,
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
func collectTests(c Collector, ch chan<- prometheus.Metric) {


	tBGP, tHTMLm, tHTMLw, bHitRateLimit, bError := c.GetTests()
	ThousandRequestsTotalMetric.Inc()

	// hint for the limit
	if bHitRateLimit {
		ThousandRequestAPILimitReached.Set(1)
	} else {
		ThousandRequestAPILimitReached.Set(0)
	}

	// if we hit the api request limit in between we skip the test details
	if bError {
		ThousandRequestsFailMetric.Inc()
		return
	}

	for e := range tBGP {

		if len(tBGP[e].Net.BgpMetrics) == 0 {
			log.Println("INFO: BGP metrics are empty for Test:", tBGP[e])
			continue
		}

		//log.Println("BGP metrics Test:", tBGP[e])
		for i := range tBGP[e].Net.BgpMetrics {

			//fmt.Println(tBGP[e].Net.Test.TestName, " | ", tBGP[e].Net.BgpMetrics[i].Prefix, " | ", tBGP[e].Net.BgpMetrics[i].MonitorName)

			// test BGP metrics
			ch <- prometheus.MustNewConstMetric(
				ThousandTestBGPReachabilityDesc,
				prometheus.GaugeValue,
				float64(tBGP[e].Net.BgpMetrics[i].Reachability),
				fmt.Sprintf("%d-%d", tBGP[e].Net.Test.TestID, i),
				tBGP[e].Net.Test.TestName,
				tBGP[e].Net.Test.Type,
				tBGP[e].Net.BgpMetrics[i].Prefix,
				tBGP[e].Net.BgpMetrics[i].CountryID,
				tBGP[e].Net.BgpMetrics[i].MonitorName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestBGPUpdatesDesc,
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
				ThousandTestBGPPathChangesDesc,
				prometheus.GaugeValue,
				float64(tBGP[e].Net.BgpMetrics[i].PathChanges),
				fmt.Sprintf("%d",tBGP[e].Net.Test.TestID),
				tBGP[e].Net.Test.TestName,
				tBGP[e].Net.Test.Type,
				tBGP[e].Net.BgpMetrics[i].Prefix,
				tBGP[e].Net.BgpMetrics[i].CountryID,
				tBGP[e].Net.BgpMetrics[i].MonitorName,
			)
		}
	}

	for e := range tHTMLm {

		if len(tHTMLm[e].Net.HTTPMetrics) == 0 {
			log.Println("INFO: HTML metrics are empty for Test:", tHTMLm[e])
			continue
		}
		for i := range tHTMLm[e].Net.HTTPMetrics {

			// test HTML metrics
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLAvgLatencyDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].AvgLatency),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLMinLatencyDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].MinLatency),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLMaxLatencyDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].MaxLatency),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLLossDesc,
				prometheus.GaugeValue,
				float64(tHTMLm[e].Net.HTTPMetrics[i].Loss),
				tHTMLm[e].Net.Test.TestName,
				tHTMLm[e].Net.Test.Type,
				tHTMLm[e].Net.Test.Prefix,
				tHTMLm[e].Net.HTTPMetrics[i].CountryID,
				tHTMLm[e].Net.HTTPMetrics[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLJitterDesc,
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
			log.Println("INFO: HTML metrics are emptry for Test:", tHTMLw[e])
			continue
		}
		for i := range tHTMLw[e].Web.HTTPServer {

			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLconnectTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ConnectTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLDNSTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].DNSTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLRedirectsDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].NumRedirects),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLreceiveTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ReceiveTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLresponseCodeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ResponseCode),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLresponseTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].ResponseTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLTotalTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].TotalTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLwaitTimeDesc,
				prometheus.GaugeValue,
				float64(tHTMLw[e].Web.HTTPServer[i].WaitTime),
				tHTMLw[e].Web.Test.TestName,
				tHTMLw[e].Web.Test.Type,
				tHTMLw[e].Web.Test.Prefix,
				tHTMLw[e].Web.HTTPServer[i].CountryID,
				tHTMLw[e].Web.HTTPServer[i].AgentName,
			)
			ch <- prometheus.MustNewConstMetric(
				ThousandTestHTMLwireSizeDesc,
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

func (t Collector) Collect(ch chan<- prometheus.Metric) {
	defer addStaticMetrics(ch)

	scrapeStart := time.Now()

	defer func() {
		if r := recover(); r != nil {
			ThousandRequestParsingFailMetric.Inc()
			log.Println("ERROR: Thousand Eyes Parsing Error (", r, ").")
		}
	}()

	collectAlerts(t, ch)
	collectTests(t, ch)


	scrapeElapsed := time.Since(scrapeStart)

	ThousandRequestScrapingTime.Add(scrapeElapsed.Seconds())

}
