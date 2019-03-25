package stores

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"

	"../structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic"
)

var news *structs.News
var err error
var client *elastic.Client
var db *sql.DB

const (
	indexName    = "app"
	docType      = "news"
	appName      = "News App"
	indexMapping = `{
		"mappings" : {
			"news" : {
				"properties" : {
					"ID" : { "type" : "integer" },
					"created" : { "type" : "date" }
				}
			}
		}
	}`
)

// Save queue to Database & ES
func Save(data []byte) {
	err = json.Unmarshal(data, &news)
	fail(err)

	toDB()
	toES()
}

// Get news from ES
func Get(page int, limit int) ([]structs.News, error) {
	ctx := context.Background()
	err = openES()
	fail(err)
	var allnews []structs.News

	var offset int
	offset = (page - 1) * limit

	sr, err := client.Search().
		Index(indexName).
		Sort("created", false).
		Type(docType).
		From(offset).
		Size(limit).
		Do(ctx)
	fail(err)

	err = openDB()
	fail(err)
	defer db.Close()

	var wg sync.WaitGroup
	var srlen = int(len(sr.Hits.Hits))
	wg.Add(srlen)

	// retrive all data on DB concurrently, by using goroutine
	for _, hit := range sr.Hits.Hits {
		go func(hit *elastic.SearchHit) {
			defer wg.Done()
			var n structs.News
			err := json.Unmarshal(*hit.Source, &n)
			fail(err)

			err = db.QueryRow("SELECT author, body FROM news where id = ?", n.ID).Scan(&n.Author, &n.Body)
			fail(err)

			allnews = append(allnews, n)
		}(hit)
	}
	wg.Wait()

	// sort result to make it ordered by Created DESC
	sort.Slice(allnews[:], func(i, j int) bool {
		return allnews[i].Created.After(allnews[j].Created)
	})

	return allnews, nil
}

func openDB() error {
	db, err = sql.Open("mysql", "root:password@tcp(db)/mydb?parseTime=true")
	if err != nil {
		return err
	}

	return nil
}

func toDB() {
	var id int64
	err = openDB()
	fail(err)
	defer db.Close()

	// insert news to mysql
	insert, err := db.Exec("INSERT INTO news (author, body) VALUES (?,?)", news.Author, news.Body)
	fail(err)

	// get last inserted record
	id, err = insert.LastInsertId()
	fail(err)
	err = db.QueryRow("SELECT id, created FROM news where id = ?", id).Scan(&news.ID, &news.Created)
	fail(err)

	fmt.Printf("Saved to DB: %d\n", id)
}

func createIndexIfNotExists(client *elastic.Client) error {
	ctx := context.Background()
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	res, err := client.CreateIndex(indexName).
		Body(indexMapping).
		Do(ctx)

	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}

	return nil
}

func openES() error {
	client, err = elastic.NewClient(elastic.SetURL("http://elastic:9200"))
	if err != nil {
		return err
	}

	err := createIndexIfNotExists(client)
	if err != nil {
		return err
	}

	return nil
}

func toES() {
	ctx := context.Background()
	err = openES()
	fail(err)

	put, err := client.Index().
		Index(indexName).
		Type(docType).
		BodyJson(news).
		Do(ctx)
	fail(err)

	fmt.Printf("Saved to ES: %s \n", put.Id)
}

func fail(err error) {
	if err != nil {
		fmt.Printf("There was an error: %s\n", err)
	}
}
