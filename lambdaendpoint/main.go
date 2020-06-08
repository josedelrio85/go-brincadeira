package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleRequest is
func HandleRequest(ctx context.Context, s3Event events.S3Event) error {

	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] ID: %s Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, record.PrincipalID.PrincipalID, s3.Bucket.Name, s3.Object.Key)
		path := strings.Split(s3.Object.Key, "/")

		for _, p := range path {
			fmt.Println(p)
		}
		principal := strings.Split(record.PrincipalID.PrincipalID, ":")

		log.Printf("Event received %s for bucket %s", time.Now().Format("2006-01-02 15-04-05"), path[1])

		endpoint := fmt.Sprintf("https://loaddatareport.bysidecar.me/import/%s/", path[1])

		log.Printf("endpoint %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Println("error lalallala")
			return err
		}
		q := req.URL.Query()
		q.Add("bucket", path[1])
		q.Add("name", path[2])
		q.Add("test", "false")
		q.Add("principal", principal[1])

		req.URL.RawQuery = q.Encode()

		log.Printf("url %s", req.URL)

		log.Println(req.URL.RawQuery)
		http := &http.Client{}
		_, err = http.Do(req)
		if err != nil {
			log.Println(err)
			log.Println("error making get request")
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

// {
//   "Records": [
//     {
//       "eventVersion": "2.0",
//       "eventSource": "aws:s3",
//       "awsRegion": "eu-west-1",
//       "eventTime": "1970-01-01T00:00:00.000Z",
//       "eventName": "ObjectCreated:Put",
//       "userIdentity": {
//         "principalId": "EXAMPLE"
//       },
//       "requestParameters": {
//         "sourceIPAddress": "127.0.0.1"
//       },
//       "responseElements": {
//         "x-amz-request-id": "EXAMPLE123456789",
//         "x-amz-id-2": "EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH"
//       },
//       "s3": {
//         "s3SchemaVersion": "1.0",
//         "configurationId": "testConfigRule",
//         "bucket": {
//           "name": "example-bucket",
//           "ownerIdentity": {
//             "principalId": "EXAMPLE"
//           },
//           "arn": "arn:aws:s3:::example-bucket"
//         },
//         "object": {
//           "key": "test/key",
//           "size": 1024,
//           "eTag": "0123456789abcdef0123456789abcdef",
//           "sequencer": "0A1B2C3D4E5F678901"
//         }
//       }
//     }
//   ]
// }
