package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/guregu/dynamo"
)

type item struct {
	ID    string `dynamo:"id,hash"`
	Title string `dynamo:"title"`
}

func (i item) HashKey() interface{} {
	return i.ID
}

func (i item) RangeKey() interface{} {
	return nil
}

type response struct {
	Message      string
	TimeMicrosec int64
}

var tableName = os.Getenv("DYNAMO_TABLE_NAME")
var port = os.Getenv("PORT")
var daxClusterURL = os.Getenv("DAX_CLUSTER_URL")
var awsRegion = os.Getenv("AWS_REGION")
var authToken = os.Getenv("AUTH_TOKEN")

var fixture = func() []interface{}{
	items := []interface{}{}
	for i := 0; i < 100; i++ {
		items = append(items, item{
			ID: strconv.Itoa(i),
			Title: strings.Repeat("A", 100),
		})
	}
	return items
}()

func newDynamoClient(tableName string, throughDax bool) (*dynamo.DB, error) {
	var (
		client dynamodbiface.DynamoDBAPI
		err    error
	)
	if throughDax {
		cfg := dax.DefaultConfig()
		cfg.HostPorts = []string{daxClusterURL}
		cfg.Region = awsRegion
		client, err = dax.New(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		client = dynamodb.New(session.New())
	}
	return dynamo.NewFromIface(client), nil
}

func auth(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := strings.Split(r.Header.Get("Authorization"), " ")
		if len(auth) < 2 || auth[1] != authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func main() {
	db, err := newDynamoClient(tableName, false)
	if err != nil {
		panic(err)
	}
	dax, err := newDynamoClient(tableName, true)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello this is ECS")
	})
	http.HandleFunc("/initialize", handleInitialize(dax))
	http.HandleFunc("/item", handleItem(db, dax))

	if port == "" {
		port = "8080"
	}
	log.Printf("Listen on :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func handleInitialize(dax *dynamo.DB) func(w http.ResponseWriter, r *http.Request) {
	return auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var items []item
			if err := dax.Table(tableName).Scan().All(&items); err != nil {
				if aerr, ok := err.(awserr.Error); !ok || aerr.Code() != dynamodb.ErrCodeResourceNotFoundException {
					log.Printf("failed to scan table: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
			var keyed = make([]dynamo.Keyed, 0, len(items))
			for _, item := range items {
				keyed = append(keyed, item)
			}
			if _, err := dax.Table(tableName).Batch("id").Write().Delete(keyed...).Run(); err != nil {
				log.Printf("failed to delete items: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if _, err := dax.Table(tableName).Batch().Write().Put(fixture...).Run(); err != nil {
				log.Printf("failed to put fixture: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, "ok")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func handleItem(db *dynamo.DB, dax *dynamo.DB) func(w http.ResponseWriter, r *http.Request) {
	return auth(func(w http.ResponseWriter, r *http.Request) {
		client := db
		if r.URL.Query().Get("dax") == "1" {
			client = dax
		}
		switch r.Method {
		case http.MethodGet:
			var item item
			before := time.Now()
			if err := client.Table(tableName).Get("id", r.URL.Query().Get("id")).One(&item); err != nil {
				if err != dynamo.ErrNotFound {
					log.Printf("failed to get an item: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
			fmt.Fprintln(w, response{
				Message:      r.URL.Query().Get("message"),
				TimeMicrosec: time.Now().Sub(before).Nanoseconds() / 1000,
			})
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			var item item
			if err := json.Unmarshal(body, &item); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			before := time.Now()
			if err := client.Table(tableName).Put(item).Run(); err != nil {
				log.Printf("failed to put an item: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, response{
				Message:      r.URL.Query().Get("message"),
				TimeMicrosec: time.Now().Sub(before).Nanoseconds() / 1000,
			})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
