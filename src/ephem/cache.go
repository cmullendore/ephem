package ephem

import (
	"log"
	"time"

	"github.com/cmullendore/ephem/src/config"
	"github.com/cmullendore/ephem/src/ephem/mysql"
)

type Cache struct {
	items  map[string]*cacheItem
	db     IDatabase
	config *config.Ephem
}

type cacheItem struct {
	key        string
	item       *string
	timestamp  time.Time
	lastAccess time.Time
	size       int
}

type sortedItems struct {
	key       string
	timestamp time.Time
	size      int
}

func NewCache(c *config.Ephem) *Cache {

	var cache Cache

	cache.items = make(map[string]*cacheItem)

	cache.config = c

	switch c.Database.Driver {
	case "mysql":
		cache.db = mysql.Initialize(c)
	}

	go cache.CleanupProcessor()

	return &cache
}

func (c *Cache) SaveItem(path *string, item *string, lifetime int, readCount int) *error {
	err := c.db.SaveItem(path, item, lifetime, readCount)
	if err != nil {
		log.Println(*err)
		var e = err
		return e
	}

	return nil
}

func (c *Cache) GetItem(path *string) (*string, *error) {

	item, err := c.db.GetItem(path)
	// This should catch zero rows returned in addition to an actual error
	if err != nil {
		return nil, err
	}

	if err := c.db.IncrementReadCount(path); err != nil {
		return nil, err
	}

	return item, nil

}

//GetItemInfo(bucket, path *string) *PersistedObjectInfo
func (c *Cache) DeleteItem(path *string) *error {
	if err := c.DeleteItem(path); err != nil {
		return err
	}

	return nil
}

func (c *Cache) CleanupProcessor() *error {
	for true {

		c.db.CleanupSecrets(c.config.MaxReads)

		dur, err := time.ParseDuration(c.config.CleanupFrequency)
		if err != nil {
			return &err
		}

		time.Sleep(dur)
	}

	return nil
}
