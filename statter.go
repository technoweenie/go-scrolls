package grohl

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Statter interface {
	Counter(sampleRate float32, bucket string, n ...int)
	Timing(sampleRate float32, bucket string, d ...time.Duration)
	Gauge(sampleRate float32, bucket string, value ...string)
}

func (c *Context) Counter(sampleRate float32, bucket string, n ...int) {
	if rand.Float32() > sampleRate {
		return
	}

	for _, num := range n {
		c.Log(Data{"metric": bucket, "count": num})
	}
}

func (c *Context) Timing(sampleRate float32, bucket string, d ...time.Duration) {
	if rand.Float32() > sampleRate {
		return
	}

	for _, dur := range d {
		c.Log(Data{"metric": bucket, "timing": int64(dur / time.Millisecond)})
	}
}

func (c *Context) Gauge(sampleRate float32, bucket string, value ...string) {
	if rand.Float32() > sampleRate {
		return
	}

	for _, v := range value {
		c.Log(Data{"metric": bucket, "gauge": v})
	}
}

type _statter struct {
	statter           Statter
	StatterSampleRate float32
	StatterBucket     string
}

func (s *_statter) SetStatter(statter Statter, sampleRate float32, bucket string) {
	if statter == nil {
		s.statter = CurrentStatter
	} else {
		s.statter = statter
	}
	s.StatterSampleRate = sampleRate
	s.StatterBucket = bucket
}

func (s *_statter) StatterBucketSuffix(suffix string) {
	sep := "."
	if strings.HasSuffix(s.StatterBucket, ".") {
		sep = ""
	}
	s.StatterBucket = s.StatterBucket + fmt.Sprintf("%s%s", sep, suffix)
}

func (s *_statter) Timing(dur time.Duration) {
	if s.statter == nil {
		return
	}

	s.statter.Timing(s.StatterSampleRate, s.StatterBucket, dur)
}

func (s *_statter) dup() *_statter {
	return &_statter{s.statter, s.StatterSampleRate, s.StatterBucket}
}
