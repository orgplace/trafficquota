package main

import (
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/orgplace/trafficquota/client"
	"github.com/orgplace/trafficquota/config"
	"github.com/orgplace/trafficquota/tokenbucket"
)

func newLogger() (*zap.Logger, error) {
	c := zap.NewProductionConfig()
	c.Level = zap.NewAtomicLevelAt(config.LogLevel)
	return c.Build()
}

func main() {
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Debug(config.Listen)
	c, err := client.NewInsecureClient(config.Listen)
	if err != nil {
		logger.Panic("failed to dial", zap.Error(err))
	}
	defer c.Close()

	testcases := []testcase{
		{
			n:              tokenbucket.DefaultBucketSize * 5,
			partitionKey:   "tenant1",
			clusteringKeys: []string{""},
		},
		{
			n:            tokenbucket.DefaultBucketSize * 5,
			partitionKey: "tenant2",
		},
		{
			n:              tokenbucket.DefaultBucketSize * 5 / 2,
			partitionKey:   "tenant3",
			clusteringKeys: []string{"operationA"},
		},
		{
			n:              tokenbucket.DefaultBucketSize * 5 / 2,
			partitionKey:   "tenant3",
			clusteringKeys: []string{"operationA", "operationB"},
		},
	}

	for i := 0; i < len(testcases); i++ {
		testcases[i].results = make(chan *result, testcases[i].n)
	}

	for i := 0; i < len(testcases); i++ {
		burst(logger, c, &testcases[i])
	}

	for i := 0; i < len(testcases); i++ {
		printResults(&testcases[i])
	}
}

type testcase struct {
	n              int
	partitionKey   string
	clusteringKeys []string
	results        chan *result
}

type result struct {
	allowed  bool
	duration time.Duration
}

func printResults(tc *testcase) {
	allow := 0
	durations := make([]time.Duration, tc.n)

	for i := 0; i < tc.n; i++ {
		r := <-tc.results
		if r.allowed {
			allow++
		}
		durations[i] = r.duration
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	fmt.Printf("partitionKey: %s, clusteringKeys: %s\n", tc.partitionKey, tc.clusteringKeys)

	fmt.Printf("allow: %d, deny: %d\n", allow, tc.n-allow)

	percentile := 1
	for i, d := range durations {
		p := float64((i+1)*10) / float64(tc.n)
		if p >= float64(percentile) {
			fmt.Printf("%6.2f%%: %s\n", p*10., d.String())
			percentile++
		}
	}
}

func burst(logger *zap.Logger, c client.Client, tc *testcase) {
	for i := 0; i < tc.n; i++ {
		go func() {
			t := time.Now()
			res, err := c.Take(tc.partitionKey, tc.clusteringKeys...)
			d := time.Since(t)
			if err != nil {
				logger.Panic("failed to take token", zap.Error(err))
			}

			tc.results <- &result{
				allowed:  res,
				duration: d,
			}
		}()
	}
}
