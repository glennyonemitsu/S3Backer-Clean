package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
)

var logger *log.Logger

// flags
var awsKey *string
var awsSecret *string
var s3Region *string
var start *int
var end *int
var width *int
var bucket *string
var prefix *string

func main() {
	if *bucket == "" {
		logger.Println("Bucket not specified")
		syscall.Exit(1)
	}

	auth := new(aws.Auth)
	auth.AccessKey = *awsKey
	auth.SecretKey = *awsSecret
	s3c := s3.New(*auth, aws.Regions[*s3Region])
	s3bucket := s3c.Bucket(*bucket)
	logger.Println(s3bucket)

	// making i with leading zeros with this format
	format := "%0" + strconv.Itoa(*width) + "d"
	for i := *start; i <= *end; i += 1 {
		suffix := fmt.Sprintf(format, i)
		key := *prefix + suffix
		logger.Printf("Deleting S3 key: %s", key)
		err := s3bucket.Del(key)
		if err != nil {
			logger.Printf("Got error deleting key: %s", err)
		}
	}
}

func init() {
	logger = log.New(os.Stdout, "", 0)

	awsKey = flag.String("awsKey", os.Getenv("AWS_ACCESS_KEY_ID"), "AWS Key. Defaults to env var AWS_ACCESS_KEY_ID")
	awsSecret = flag.String("awsSecret", os.Getenv("AWS_SECRET_KEY"), "AWS Secret. Defaults to env var AWS_SECRET_KEY")
	s3Region = flag.String("s3Region", "us-east-1", "AWS S3 region")
	start = flag.Int("start", 0, "Starting number")
	end = flag.Int("end", 0, "Ending number (inclusive)")
	width = flag.Int("width", 6, "Key number width (ex. when width = 6, 1 = 000001)")
	bucket = flag.String("bucket", "", "Bucket")
	prefix = flag.String("prefix", "/", "Key prefix")
	flag.Parse()
}

