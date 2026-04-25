package s3test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Result struct {
	PutLatency time.Duration
	GetLatency time.Duration
	Size       int
	Error      error
}

func Run(ctx context.Context, client *s3.Client, bucket, key string, payload []byte) Result {
	var res Result
	res.Size = len(payload)

	start := time.Now()
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(payload),
	})
	res.PutLatency = time.Since(start)

	if err != nil {
		log.Printf("Error putting object: %v", err)
		res.Error = fmt.Errorf("put error: %w", err)
		return res
	}

	start = time.Now()
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	res.GetLatency = time.Since(start)

	if err != nil {
		log.Printf("Error getting object: %v", err)
		res.Error = fmt.Errorf("get error: %w", err)
		return res
	}
	defer out.Body.Close()

	return res
}
