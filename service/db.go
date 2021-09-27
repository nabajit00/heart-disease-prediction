package service

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Stats_t struct {
	UniqueUsers    int
	TotalPred      int
	PredictedCHD   int
	PredictedNoCHD int
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db/data.db")
	checkErr(err)

	statement, err := db.Prepare("create table if not exists Ip(ip char(45) Primary Key)")
	checkErr(err)
	statement.Exec()

	statement, err = db.Prepare("create table if not exists Prediction(ip char(45),pred INT8,foreign key(ip) references Ip(ip))")
	checkErr(err)
	statement.Exec()
}

func Test() {
	statement, err := db.Prepare("insert into Ip values(?)")
	checkErr(err)
	res, err := statement.Exec("123.123")
	checkErr(err)
	fmt.Println(res)
}

func InsertResult(ip string, p int8) error {
	row := db.QueryRow("Select Count(*) from Ip where ip=?", ip)

	var count int
	err := row.Scan(&count)
	checkErr(err)

	if err != nil {
		return err
	}
	if count == 0 {
		statement, err := db.Prepare("insert into Ip values(?)")
		if err != nil {
			return err
		}
		_, err = statement.Exec(ip)
		if err != nil {
			return err
		}
	}

	statement, err := db.Prepare("insert into Prediction values(?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(ip, p)
	if err != nil {
		return err
	}
	return nil
}

func GetStats() (Stats_t, error) {
	var stats Stats_t
	s := db.QueryRow("select count(*) from Ip")
	if s.Err() != nil {
		return stats, s.Err()
	}
	s.Scan(&(stats.UniqueUsers))

	s = db.QueryRow("select Count(pred) from Prediction")
	if s.Err() != nil {
		return stats, s.Err()
	}
	s.Scan(&(stats.TotalPred))

	s = db.QueryRow("select Count(pred) from Prediction where pred=1")
	if s.Err() != nil {
		return stats, s.Err()
	}
	s.Scan(&(stats.PredictedCHD))

	s = db.QueryRow("select Count(pred) from Prediction where pred=0")
	if s.Err() != nil {
		return stats, s.Err()
	}
	s.Scan(&(stats.PredictedNoCHD))

	return stats, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
