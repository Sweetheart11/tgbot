package fetcher

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Sweetheart11/tgbot/lib/dstr"
	"github.com/Sweetheart11/tgbot/model"
	"github.com/Sweetheart11/tgbot/source"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}

type SourceStorage interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles ArticleStorage
	sources  SourceStorage

	fetchInterval  time.Duration
	filterKeywords []string
}

func New(
	articleStorage ArticleStorage,
	sourceStorage SourceStorage,
	fetchInterval time.Duration,
	filterKeywords []string,
) *Fetcher {
	return &Fetcher{
		articles:       articleStorage,
		sources:        sourceStorage,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, s := range sources {
		wg.Add(1)

		rssSource := source.NewRSSSourceFromModel(s)
		go func(source Source) error {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("error fetching items from source: %s: %v", source.Name(), err)
				return err
			}

			if err := f.processItems(ctx, source, items); err != nil {
				log.Printf("error processing items from source: %s: %v", source.Name(), err)
				return err
			}

			return nil
		}(rssSource)
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemShouldBeSkipped(item) {
			continue
		}

		if err := f.articles.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Summary:     item.Summary,
			Link:        item.Link,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool {
	categoriesSet := dstr.New(item.Categories...)
	for _, keyword := range f.filterKeywords {
		titleContainsKeyword := strings.Contains(strings.ToLower(item.Title), keyword)
		if categoriesSet.Contains(keyword) || titleContainsKeyword {
			return true
		}
	}

	return false
}
