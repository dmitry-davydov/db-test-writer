package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type MysqlTest struct {
	dbConnection *sql.DB
}

func NewMysqlTest(dsn string) *MysqlTest {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return &MysqlTest{dbConnection: db}
}

func (t MysqlTest) Write(stat DBStat) error {

	_, err := t.dbConnection.Exec("INSERT INTO `stats` (`referer`, `useragent`, `ip`, `country_id`, `operator_id`, `land_id`, `partner_id`, `browser_id`, `platform_id`, `model_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", stat.Referer,
		stat.Useragent,
		stat.Ip,
		stat.CountryId,
		stat.OperatorId,
		stat.LandId,
		stat.PartnerId,
		stat.Browser,
		stat.PlatformId,
		stat.ModelId)
	if err != nil {
		panic(err)
	}
	return err
}
