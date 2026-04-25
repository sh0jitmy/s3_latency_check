package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"s3-latency/internal/s3test"
	"s3-latency/internal/stats"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

type nopWarnLogger struct{}

func (l nopWarnLogger) Logf(classification logging.Classification, format string, v ...interface{}) {

	// WARNを無視
	if classification == logging.Warn {
		return
	}
	log.Printf(format, v...)

}

func main() {
	ctx := context.Background()

	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithLogger(nopWarnLogger{}),
	)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		o.Region = os.Getenv("AWS_REGION")
	})

	iterations := 100

	batch := s3test.RunBatch(ctx, client, bucket, iterations)

	putStats := stats.Calculate(batch.PutLatencies)
	getStats := stats.Calculate(batch.GetLatencies)

	errorRate := float64(batch.Errors) / float64(iterations) * 100

	fmt.Printf("PUT: %+v\n", putStats)
	fmt.Printf("GET: %+v\n", getStats)
	fmt.Printf("ErrorRate: %.2f%%\n", errorRate)

	fmt.Printf(
		"SLACK_MESSAGE=PUT[p95:%v p99:%v] GET[p95:%v p99:%v] Err:%.2f%%",
		putStats.P95,
		putStats.P99,
		getStats.P95,
		getStats.P99,
		errorRate,
	)
}
