package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

func makePostStruct(item RSSItem, feedId uuid.UUID) database.CreatePostParams {
	t, err := time.Parse(time.RFC1123Z, item.PubDate)
	if err != nil {
		return database.CreatePostParams{}
	}
	str_is_val := item.Description != ""
	return database.CreatePostParams{
		ID:        generateUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     item.Title,
		Url:       item.Link,
		Description: sql.NullString{
			String: item.Description,
			Valid:  str_is_val,
		},
		PublishedAt: sql.NullTime{
			Time:  t,
			Valid: true,
		},
		FeedID: feedId,
	}
}

func (apiConf *apiConfig) startWorker(num_items, delayInSeconds int) {
	log.Println("BG_WORKER_LOG: starting worker")
	wg := sync.WaitGroup{}
	for {
		// get URLS from DB
		log.Println("BG_WORKER_LOG: getting feeds to fetch")
		feeds, err := apiConf.DB.GetNextFeedsToFetch(context.Background(), int32(num_items))
		if err != nil {
			continue
		}

		// fetch those URLS using fetch func below
		log.Println("BG_WORKER_LOG: fetching feeds. . .")
		for _, feed := range feeds {
			wg.Add(1)
			go func(feed database.Feed) {
				defer wg.Done()
				RssBody, err := fetchData(feed.Url)
				if err != nil {
					log.Printf("Error while fetching %s\n", feed.Url)
					return
				}
				// todo save this data to the DB.
				for _, item := range RssBody.Channel.Item {
					apiConf.DB.CreatePost(context.Background(), makePostStruct(item, feed.ID))
				}
			}(feed)
			apiConf.DB.MarkFeedFetched(context.Background(),
				database.MarkFeedFetchedParams{
					ID: feed.ID, UpdatedAt: time.Now(),
					LastFetchedAt: sql.NullTime{Time: time.Now(),
						Valid: true,
					}})
		}

		// wait
		wg.Wait()
		log.Println("BG_WORKER_LOG: done fetching. delaying for next cycle.")
		time.Sleep(time.Second * time.Duration(delayInSeconds))
	}
}

type RSSItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type RssBody struct {
	Channel struct {
		Title string `xml:"title"`
		Link  struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func fetchData(url string) (RssBody, error) {
	rBody := RssBody{}

	resp, err := http.Get(url)
	if err != nil {
		return RssBody{}, err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return RssBody{}, err
	}

	err = xml.Unmarshal(body, &rBody)
	resp.Body.Close()

	if err != nil {
		return RssBody{}, err
	}
	return rBody, nil
}
