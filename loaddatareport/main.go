package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	log.Printf("Thisvaina process launched at %s", time.Now().Format("2006-01-02 15-04-05"))

	filename := "2020-06-03.sql.gz"

	if err := getFromS3(filename); err != nil {
		log.Fatal(err)
	}

	if err := unzip(filename); err != nil {
		log.Fatal(err)
	}
	newfilename := strings.TrimSuffix(filename, ".gz")

	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := struct {
		Host   string
		User   string
		Pass   string
		Port   int64
		Dbname string
		File   string
	}{
		Host:   getSetting("DB_HOST"),
		User:   getSetting("DB_USER"),
		Pass:   getSetting("DB_PASS"),
		Dbname: "testing",
		Port:   portInt,
		// Host:   "leads-pre.c848y92oajny.eu-west-1.rds.amazonaws.com",
		// User:   "leads",
		// Pass:   "LW3PBzuqy3zfBrqBbbFM",
		// Dbname: "leads",
		// File:   "select * from leads order by id desc limit 10;",
		File: newfilename,
	}

	host := fmt.Sprintf("-h%s", database.Host)
	user := fmt.Sprintf("-u%s", database.User)
	pass := fmt.Sprintf("-p%s", database.Pass)
	portt := fmt.Sprintf("-P %d", database.Port)
	db := fmt.Sprintf("-D%s", database.Dbname)
	// file := fmt.Sprintf("-e %s", database.File)
	file := fmt.Sprintf("source ./%s", database.File)

	log.Println(host)
	log.Println(user)
	log.Println(pass)
	log.Println(portt)
	log.Println(db)
	log.Println(file)
	// cmd := exec.Command("/usr/bin/mysql", "-h127.0.0.1", "-P 3306", "-uroot", "-proot_bsc", "-f", "-Dwebservice", "-e show tables;")
	// cmd := exec.Command("/usr/bin/mysql", host, portt, user, pass, db, "-f", file)
	// cmd := exec.Command("/usr/bin/mysql", "-h127.0.0.1", "-P 3306", "-uroot", "-proot_bsc", "-f", "-Dwebservice", "-e ", "source filename.sql")
	cmd := exec.Command("/usr/bin/mysql", host, portt, user, pass, db, "-f", "-e", file)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bytes))

	err = ioutil.WriteFile("./out.sql", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	byteserr, err := ioutil.ReadAll(stderr)
	err = ioutil.WriteFile("./err.sql", byteserr, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(byteserr))

	if err := os.Remove(filename); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s removed succesfully", filename)

	if err := os.Remove(newfilename); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s removed succesfully", newfilename)

	log.Printf("Thisvaina process ended at %s", time.Now().Format("2006-01-02 15-04-05"))
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}

func getFromS3(filename string) error {
	// The session the S3 Downloader will use
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	mybucket := "data.bysidecar.me/backups/ws/"

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
		Bucket: aws.String(mybucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return fmt.Errorf("Unable to download item %q, %v", filename, err)
	}
	log.Println("Downloaded", f.Name(), numBytes, "bytes")
	return nil
}

func listBuckets() {
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

func unzip(filename string) error {
	gzipfile, err := os.Open(filename)

	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		return err
	}
	defer reader.Close()

	newfilename := strings.TrimSuffix(filename, ".gz")

	writer, err := os.Create(newfilename)

	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	log.Printf("%s unzipped succesful", newfilename)
	return nil
}
