package badgerkv

import (
	"sync"

	badger "github.com/dgraph-io/badger/v4"
)

var (
	badgerDB   *BadgerDB
	badgerOnce sync.Once
	c          *Config = &Config{
		PersistDir: "./pv",
	}
)

type BadgerDB struct {
	db *badger.DB
}

func NewBadgerDB(file string) *BadgerDB {
	db, err := badger.Open(badger.DefaultOptions(file))
	if err != nil {
		panic(err)
	}

	return &BadgerDB{
		db: db,
	}
}

type Config struct {
	PersistDir string
}

func SetConfig(conf *Config) {
	c = conf
}

func Init() {
	badgerOnce.Do(func() {
		badgerDB = NewBadgerDB(c.PersistDir)
	})
}

func GetBadgerDB() *BadgerDB {
	if badgerDB == nil {
		Init()
	}

	return badgerDB
}

func (db *BadgerDB) GetRawDB() *badger.DB {
	return db.db
}

func (db *BadgerDB) Set(k, v string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(k), []byte(v))
	})
}

func (db *BadgerDB) Get(k string) (string, error) {
	var v []byte
	if err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}

		if v, err = item.ValueCopy(nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	return string(v), nil
}

func (db *BadgerDB) Delete(k string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(k))
	})
}

func (db *BadgerDB) Range() (map[string]string, error) {
	res := map[string]string{}

	if err := db.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			k := item.Key()
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			res[string(k)] = string(v)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}
