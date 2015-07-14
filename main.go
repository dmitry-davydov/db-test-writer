package main

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type DBTester interface {
	Write(DBStat) error
}

type DBStat struct {
	Id         int32  `json:"id"`
	Referer    string `json:"referer"`
	Useragent  string `json:"useragent"`
	Ip         string `json:"ip"`
	CountryId  int32  `json:"countryId"`
	OperatorId int32  `json:"operatorId"`
	LandId     int32  `json:"landId"`
	PartnerId  int32  `json:"partnerId"`
	Browser    int32  `json:"browser"`
	PlatformId int32  `json:"platformId"`
	ModelId    int32  `json:"modelId"`
}

type data struct {
	name  string
	count int
}

type statistics struct {
	sync.RWMutex
	stat   map[string]int
	ticker *time.Ticker
}

func (s *statistics) printStat(t time.Time) {
	s.Lock()

	totalWriteCount = totalWriteCount + s.stat["write"]

	fmt.Printf("|%20s|%20d|%20d|\n", fmt.Sprintf("%2d:%2d:%2d", t.Hour(), t.Minute(), t.Second()), s.stat["write"], totalWriteCount)

	for k, _ := range s.stat {
		s.stat[k] = 0
	}

	s.Unlock()
}

func (s *statistics) captureStat(dataItem data) {
	if _, ok := s.stat[dataItem.name]; ok {
		s.stat[dataItem.name] = s.stat[dataItem.name] + dataItem.count
	} else {
		s.stat[dataItem.name] = dataItem.count
	}
}

func newStatistics() *statistics {
	return &statistics{
		stat:   map[string]int{"read": 0, "write": 0},
		ticker: time.NewTicker(time.Second)}
}

var totalWriteCount int
var driver string
var dsn string
var concurentDbWrite int

func init() {
	flag.StringVar(&driver, "driver", "", "mysql|mongo|couchbase|elastic")
	flag.StringVar(&dsn, "dsn", "", "user:password@/dbname mongodb://localhost/stat couchbase://127.0.0.1:8091 127.0.0.1:9200")
	flag.IntVar(&concurentDbWrite, "c", 20, "Number of concurent db writing, default 20")
}

func main() {

	flag.Parse()

	totalWriteCount = 0
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)

	st := newStatistics()

	var test DBTester

	switch driver {
	case "mysql":
		test = NewMysqlTest(dsn)

	case "mongo":
		test = NewMongoTest(dsn)
	case "couchbase":
		test = NewCouchBaseTest(dsn)
	case "elastic":
		test = NewElasticTest(dsn)
	default:
		panic("You should choose driver")
	}

	fmt.Println("GOMAXPROCS ", cpuNum)
	fmt.Println("Starting banchmark")
	fmt.Printf("|%20s|%20s|%20s|\n", "Time", "Write per sec", "Total write count")

	for i := 0; i < concurentDbWrite; i++ {
		go writeToDB(st, test)
	}

	for t := range st.ticker.C {
		go st.printStat(t)
	}

}

func writeToDB(st *statistics, test DBTester) {
	for i := 0; i < 1000000; i++ {
		dbData := &DBStat{Id: 0,
			Referer:    randStr(20, "alpha"),
			Useragent:  randStr(20, "alpha"),
			Ip:         randStr(20, "alpha"),
			CountryId:  rand.Int31n(20),
			OperatorId: rand.Int31n(20),
			LandId:     rand.Int31n(20),
			PartnerId:  rand.Int31n(20),
			Browser:    rand.Int31n(20),
			PlatformId: rand.Int31n(20),
			ModelId:    rand.Int31n(20)}

		if err := test.Write(*dbData); err == nil {
			dataItem := &data{name: "write", count: 1}

			st.captureStat(*dataItem)
		}
	}
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	crand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
