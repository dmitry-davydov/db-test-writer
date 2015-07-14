package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoDbStruct struct {
	Id         bson.ObjectId `bson:"_id"`
	Referrer   string        `bson:"referrer"`
	Useragent  string        `bson:"useragent"`
	IP         string        `bsob:"ip"`
	CountryId  int32         `bson:"countryId"`
	OperatorId int32         `bson:"operatorId"`
	LandId     int32         `bson:"landId"`
	PartnerId  int32         `bson:"partnerId"`
	Browser    int32         `bson:"browserId"`
	PlatformId int32         `bson:"platformId"`
	ModelId    int32         `bson:"modelId"`
}

type MongoTest struct {
	dbConnection *mgo.Session
}

func NewMongoTest(dsn string) *MongoTest {
	db, err := mgo.Dial(dsn)
	if err != nil {
		panic(err)
	}

	return &MongoTest{dbConnection: db}

}

func (t MongoTest) Write(stat DBStat) error {
	collection := t.dbConnection.DB("stat").C("stat")

	insertDocument := &mongoDbStruct{Id: bson.NewObjectId(),
		Referrer:   stat.Referer,
		Useragent:  stat.Useragent,
		IP:         stat.Ip,
		CountryId:  stat.CountryId,
		OperatorId: stat.OperatorId,
		LandId:     stat.LandId,
		PartnerId:  stat.PartnerId,
		Browser:    stat.Browser,
		PlatformId: stat.PlatformId,
		ModelId:    stat.ModelId}

	err := collection.Insert(insertDocument)

	return err
}
