package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
)

var logger *log.Logger

// flags
var awsKey string
var awsSecret string
var s3Region string
var start int
var end int
var width int
var bucket string
var prefix string

func main() {
	if bucket == "" {
		logger.Println("Bucket not specified")
		syscall.Exit(1)
	}

	auth := new(aws.Auth)
	auth.AccessKey = awsKey
	auth.SecretKey = awsSecret
	s3c := s3.New(*auth, aws.Regions[s3Region])
	s3bucket := s3c.Bucket(bucket)

	// making i with leading zeros with this format
	format := fmt.Sprintf("%%0%dd", width)
	for i := start; i <= end; i += 1 {
		suffix := fmt.Sprintf(format, i)
		key := prefix + suffix
		logger.Printf("Deleting S3 key: %s%s", bucket, key)
		err := s3bucket.Del(key)
		if err != nil {
			logger.Printf("Got error deleting key: %s", err)
		}
	}
}

func init() {
	logger = log.New(os.Stdout, "", 0)

	flag.StringVar(&awsKey, "awsKey", os.Getenv("AWS_ACCESS_KEY_ID"), "AWS Key. Defaults to env var AWS_ACCESS_KEY_ID")
	flag.StringVar(&awsSecret, "awsSecret", os.Getenv("AWS_SECRET_KEY"), "AWS Secret. Defaults to env var AWS_SECRET_KEY")
	flag.StringVar(&s3Region, "s3Region", "us-east-1", "AWS S3 region")
	flag.IntVar(&start, "start", 0, "Starting number")
	flag.IntVar(&end, "end", 0, "Ending number (inclusive)")
	flag.IntVar(&width, "width", 6, "Key number width (ex. when width = 6, 1 = 000001)")
	flag.StringVar(&bucket, "bucket", "", "Bucket")
	flag.StringVar(&prefix, "prefix", "/", "Key prefix")
	flag.Parse()
}
