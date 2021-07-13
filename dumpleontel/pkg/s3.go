package dumpfootel

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// GetFromS3 retrieve a filename from S3 Bucket
func GetFromS3(filename string, bucket string) error {
	log.Printf("Downloading process started at %s", time.Now().Format("2006-01-02 15-04-05"))

	// The session the S3 Downloader will use
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	bucketpath := "data.josedelrio85.me/backups/"
	bucket = fmt.Sprintf("%s%s/", bucketpath, bucket)
	log.Printf("bucket: %s", bucket)
	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}
	defer f.Close()

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Write the contents of S3 Object to the file
	numBytes, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return fmt.Errorf("Unable to download item %q, %v", filename, err)
	}

	log.Println("Downloaded", f.Name(), numBytes, "bytes")
	log.Printf("Downloading process finished at %s", time.Now().Format("2006-01-02 15-04-05"))
	return nil
}

// ListBuckets list buckets from S3
func ListBuckets() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, bucket := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(bucket.Name), aws.TimeValue(bucket.CreationDate))

		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(*bucket.Name)})
		if err != nil {
			exitErrorf("Unable to list items in bucket %q, %v", *bucket.Name, err)
		}

		for _, item := range resp.Contents {
			fmt.Println("Name:         ", *item.Key)
			fmt.Println("Last modified:", *item.LastModified)
			fmt.Println("Size:         ", *item.Size)
			fmt.Println("Storage class:", *item.StorageClass)
			fmt.Println("")
		}
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
