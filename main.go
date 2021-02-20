package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type result_struct struct {
	Records []struct {
		EventVersion string `json:"eventVersion"`
		EventSource  string `json:"eventSource"`
		AwsRegion    string `json:"awsRegion"`
		EventTime    string `json:"eventTime"`
		EventName    string `json:"eventName"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			XAmzRequestID string `json:"x-amz-request-id"`
			XAmzID2       string `json:"x-amz-id-2"`
		} `json:"responseElements"`
		S3 struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationID string `json:"configurationId"`
			Bucket          struct {
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
				Arn string `json:"arn"`
				ID  string `json:"id"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int    `json:"size"`
				Etag      string `json:"etag"`
				VersionID string `json:"versionId"`
				Sequencer string `json:"sequencer"`
				Metadata  []struct {
					Key string `json:"key"`
					Val string `json:"val"`
				} `json:"metadata"`
				Tags []interface{} `json:"tags"`
			} `json:"object"`
		} `json:"s3"`
		EventID    string `json:"eventId"`
		OpaqueData string `json:"opaqueData"`
	} `json:"Records"`
}

var s3_objectCreated_wildcard = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectCreated_wildcard",
		Help: "Total number of requests to root",
	},
)
var s3_objectCreated_put = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectCreated_put",
		Help: "Total number of requests to root",
	},
)
var s3_objectCreated_post = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectCreated_post",
		Help: "Total number of requests to root",
	},
)
var s3_objectCreated_copy = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectCreated_copy",
		Help: "Total number of requests to root",
	},
)
var s3_objectCreated_completeMultipartUpload = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectCreated_completeMultipartUpload",
		Help: "Total number of requests to root",
	},
)
var s3_objectRemoved_wildcard = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectRemoved_wildcard",
		Help: "Total number of requests to root",
	},
)
var s3_objectRemoved_delete = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectRemoved_delete",
		Help: "Total number of requests to root",
	},
)
var s3_objectRemoved_deleteMarkerCreated = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "s3_objectRemoved_deleteMarkerCreated",
		Help: "Total number of requests to root",
	},
)

func body(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	//log.Println(string(body))
	var t result_struct
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println("New Operation:", t.Records[0].EventName)
	if t.Records[0].EventName == "s3:ObjectCreated:*" {
		s3_objectCreated_wildcard.Inc()
	} else if t.Records[0].EventName == "s3:ObjectCreated:Put" {
		s3_objectCreated_put.Inc()
	} else if t.Records[0].EventName == "s3:ObjectCreated:Post" {
		s3_objectCreated_post.Inc()
	} else if t.Records[0].EventName == "s3:ObjectCreated:Copy" {
		s3_objectCreated_copy.Inc()
	} else if t.Records[0].EventName == "s3:ObjectCreated:CompleteMultipartUpload" {
		s3_objectCreated_completeMultipartUpload.Inc()
	} else if t.Records[0].EventName == "s3:ObjectRemoved:*" {
		s3_objectRemoved_wildcard.Inc()
	} else if t.Records[0].EventName == "s3:ObjectRemoved:Delete" {
		s3_objectRemoved_delete.Inc()
	} else if t.Records[0].EventName == "s3:ObjectRemoved:DeleteMarkerCreated" {
		s3_objectRemoved_deleteMarkerCreated.Inc()
	} else {
		log.Println("Event not supported:", t.Records[0].EventName)
	}
}

func main() {
	http.HandleFunc("/test", body)
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(s3_objectCreated_wildcard)
	prometheus.MustRegister(s3_objectCreated_put)
	prometheus.MustRegister(s3_objectCreated_post)
	prometheus.MustRegister(s3_objectCreated_copy)
	prometheus.MustRegister(s3_objectCreated_completeMultipartUpload)
	prometheus.MustRegister(s3_objectRemoved_wildcard)
	prometheus.MustRegister(s3_objectRemoved_delete)
	prometheus.MustRegister(s3_objectRemoved_deleteMarkerCreated)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
