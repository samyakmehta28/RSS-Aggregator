package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samyakmehta28/RSS-Aggregator/internal/database"
)

func startScraper(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	fmt.Printf("Starting scraper with %d workers\n", concurrency)
	ticker := time.NewTicker(timeBetweenRequest)
	for ;; <-ticker.C{
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			fmt.Printf("Error fetching feeds: %v\n", err)
			continue
		}
		
		wg := &sync.WaitGroup{}
		for i,_ := range feeds {
			wg.Add(1)
			go func (feed database.Feed) {
				defer wg.Done()
				_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
				if err != nil {
					fmt.Printf("Error marking feed as fetched: %v\n", err)
					return
				}
				rssFeed, err := urlToRSS(feed.Url)
				if err != nil {
					fmt.Printf("Error fetching RSS Fed: %v\n", err)
					return
				}
				

				for _, item := range rssFeed.Channel.Item {
					if item.Title == "" {
						continue
					}

					_, err := db.CreatePost(context.Background(), database.CreatePostParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Title:       item.Title,
						Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
						PublishedAt: func() time.Time {
							parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
							if err != nil {
								fmt.Printf("Error parsing publication date: %v\n", err)
								return time.Now()
							}
							return parsedTime
						}(),
						Url:         item.Link,
						FeedID:      feed.ID,
					})

					if err != nil {

						if strings.Contains(err.Error(), "duplicate key ") {
							continue
						}
						fmt.Printf("Error creating post: %v\n", err)
						continue
					}
				}

				log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
			}(feeds[i])
		}
		wg.Wait()
	}
}