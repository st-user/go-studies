package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	mux.HandleFunc("/getObject", func(w http.ResponseWriter, r *http.Request) {

		// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
		// https://docs.aws.amazon.com/AmazonS3/latest/userguide/example-policies-s3.html
		// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))

		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Error 1 !!!"))
			return
		}

		client := s3.NewFromConfig(cfg)
		out, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String("gs-test-bucket"),
			Key:    aws.String("testdata.txt"),
		})

		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Error 2 !!!"))
			return
		}

		defer out.Body.Close()

		data, err := io.ReadAll(out.Body)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Error 3 !!!"))
			return
		}
		fmt.Fprintf(w, "Got %s\n", string(data))
	})

	port := 80
	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	fmt.Printf("Listening on %d\n", port)
	log.Fatal(server.ListenAndServe())
}
