package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

type S3Bucket struct {
	// Date the bucket was created.
	CreationDate time.Time `json:"created_on"`

	// The name of the bucket.
	Name string `json:"name"`
}

type S3Buckets struct {
	Region string
}

type S3Objects struct {
	Region     string
	Bucketname string
}

type S3ObjectsPrefix struct {
	Region     string
	Bucketname string
	Prefix     string
	Assetid    string
}

// list objects in a bucket
func listObjects(bucket string, prefix string) error {

	// Create S3 service client
	svc := s3.New(Sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for i, item := range resp.Contents {
		if i < 10 {
			fmt.Println("Name:         ", *item.Key)
			fmt.Println("Last modified:", *item.LastModified)
			fmt.Println("Size:         ", *item.Size)
			fmt.Println("Storage class:", *item.StorageClass)
			fmt.Println("")
		}
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")

	return nil
}

// listBuckets lists the S3 buckets with the given prefix e.g. "video-sample"
func listBucketsPrefix(prefix string) error {
	// Create S3 service client
	svc := s3.New(Sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Println("Could not list buckets")
		return err
	}

	for _, b := range result.Buckets {
		if strings.HasPrefix(*b.Name, prefix) {
			fmt.Println(*b.Name)
		}
	}

	return nil
}

// listBuckets lists the S3 buckets
//func listBuckets() ([]S3Bucket, error) {
func listBuckets() error {

	// S3 service client
	svc := s3.New(Sess)

	//result, err := svc.ListBuckets(nil)
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	var buckets []S3Bucket
	for _, b := range result.Buckets {
		buckets = append(buckets, S3Bucket{Name: *(b.Name),
			CreationDate: *(b.CreationDate)})
	}

	// TODO: comment out - For debug only
	for i, b := range buckets {
		if i < 10 {
			fmt.Printf("* %s created on %s\n",
				b.Name, b.CreationDate)
		}
	}
	//	return buckets, nil
	return nil

}
