package s3test

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BatchResult struct {
	PutLatencies []time.Duration
	GetLatencies []time.Duration
	Errors       int
}

func RunBatch(ctx context.Context, client *s3.Client, bucket string, iterations int) BatchResult {
	var result BatchResult

	payload := make([]byte, 1024*100)
	rand.Read(payload)

	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("latency-%d-%d", time.Now().UnixNano(), i)

		res := Run(ctx, client, bucket, key, payload)

		if res.Error != nil {
			result.Errors++
			continue
		}

		result.PutLatencies = append(result.PutLatencies, res.PutLatency)
		result.GetLatencies = append(result.GetLatencies, res.GetLatency)
	}

	return result
}
