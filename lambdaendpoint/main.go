package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleRequest is
func HandleRequest(ctx context.Context, s3Event events.S3Event) error {

	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
		path := strings.Split(s3.Object.Key, "/")

		for _, p := range path {
			fmt.Println(p)
		}

		payload := struct {
			Bucket string `json:"bucket"`
			Name   string `json:"name"`
			Test   bool   `json:"test"`
		}{
			Bucket: path[1],
			Name:   path[2],
		}

		endpoint := fmt.Sprintf("https://algo.bysidecar.me/import/%s", payload.Bucket)

		bytevalues, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, _ := ioutil.ReadAll(resp.Body)
		response := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{}
		if err := json.Unmarshal(data, &response); err != nil {
			return err
		}
		log.Println(response)
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
