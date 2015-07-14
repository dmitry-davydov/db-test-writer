package main

import (
	"github.com/couchbase/go-couchbase"
	"github.com/satori/go.uuid"
)

type CouchBaseTest struct {
	dbConnection couchbase.Client
	bucket       *couchbase.Bucket
}

func NewCouchBaseTest(dsn string) *CouchBaseTest {
	db, err := couchbase.Connect(dsn)
	if err != nil {
		panic(err)
	}

	pool, err := db.GetPool("default")
	if err != nil {
		panic(err)
	}

	bucket, err := pool.GetBucket("default")
	if err != nil {
		panic(err)
	}

	return &CouchBaseTest{dbConnection: db, bucket: bucket}

}

func (t CouchBaseTest) Write(stat DBStat) error {

	key := uuid.NewV4()
	err := t.bucket.Set(key.String(), 0, stat)

	return err
}
