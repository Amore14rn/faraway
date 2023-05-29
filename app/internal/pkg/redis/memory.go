package redis

import (
	"sync"
	"time"
)

type TimeKeeper interface {
	Current() time.Time
}

type TemporalStorage struct {
	cacheStore map[int]temporalData
	mutualExcl *sync.Mutex
	timeSystem TimeKeeper
}

type temporalData struct {
	TimeStamp int64
	LifeSpan  int64
}

func InitializeTemporalStorage(timeSystem TimeKeeper) *TemporalStorage {
	return &TemporalStorage{
		cacheStore: make(map[int]temporalData, 0),
		mutualExcl: &sync.Mutex{},
		timeSystem: timeSystem,
	}
}

func (store *TemporalStorage) Insert(key int, lifeSpan int64) error {
	store.mutualExcl.Lock()
	defer store.mutualExcl.Unlock()
	store.cacheStore[key] = temporalData{
		TimeStamp: store.timeSystem.Current().Unix(),
		LifeSpan:  lifeSpan,
	}
	return nil
}

func (store *TemporalStorage) Retrieve(key int) (bool, error) {
	store.mutualExcl.Lock()
	defer store.mutualExcl.Unlock()
	data, present := store.cacheStore[key]
	if present && store.timeSystem.Current().Unix()-data.TimeStamp > data.LifeSpan {
		return false, nil
	}
	return present, nil
}

func (store *TemporalStorage) Remove(key int) {
	store.mutualExcl.Lock()
	defer store.mutualExcl.Unlock()
	delete(store.cacheStore, key)
}
