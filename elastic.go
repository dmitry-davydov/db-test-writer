package main

import (
	"github.com/belogik/goes"
	"github.com/satori/go.uuid"
	"net/url"
	"strings"
)

type ElasticTest struct {
	dbConnection *goes.Connection
}

func NewElasticTest(dsn string) *ElasticTest {

	splitted := strings.Split(dsn, ":")
	host := splitted[0]
	port := splitted[1]

	db := goes.NewConnection(host, port)

	return &ElasticTest{dbConnection: db}

}

func (t ElasticTest) Write(stat DBStat) error {

	key := uuid.NewV4()

	document := goes.Document{
		Index:  "stat",
		Type:   "stat",
		Id:     key.String(),
		Fields: stat,
	}

	extraArgs := make(url.Values, 1)
	extraArgs.Set("ttl", "86400000")

	_, err := t.dbConnection.Index(document, extraArgs)

	return err
}
