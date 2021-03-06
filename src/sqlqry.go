package main

import (
	"log"
	"time"
	"io/ioutil"
	"encoding/json"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type sqlConf struct {
	Port string
	IP string
	Username string
	Password string
}

var SQLConn sqlConf

func initSQL() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	err = json.Unmarshal(b, &SQLConn)
	if err != nil {
		log.Println("ERROR:", err)
	}
}

func checkAccount(user, pass string) bool {
	var res int

	db, err := sql.Open("mysql",
		SQLConn.Username+":"+SQLConn.Password+"@tcp("+SQLConn.IP+":"+SQLConn.Port+")/plugngo")
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}
	defer db.Close()
	rows, err := db.Query("SELECT COUNT(*) FROM admin WHERE user=? AND password=?",
		user, pass)
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&res)
	if err != nil || res == 0 {
		return false
	}
	return true
}

func addAdmin(user, pass string) error {
	db, err := sql.Open("mysql",
		SQLConn.Username+":"+SQLConn.Password+"@tcp("+SQLConn.IP+":"+SQLConn.Port+")/plugngo")
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO admin (user, password) VALUE (?, ?)",
		user, pass)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	return nil
}

func hasAdmin() bool {
	var res int

	db, err := sql.Open("mysql",
		SQLConn.Username+":"+SQLConn.Password+"@tcp("+SQLConn.IP+":"+SQLConn.Port+")/plugngo")
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}
	defer db.Close()
	rows, err := db.Query("SELECT COUNT(*) FROM admin")
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&res)
	if err != nil || res == 0 {
		return false
	}
	return true
}

func storeDatum(id string, power uint32, voltage byte) error {
	db, err := sql.Open("mysql",
		SQLConn.Username+":"+SQLConn.Password+"@tcp("+SQLConn.IP+":"+SQLConn.Port+")/plugngo")
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO status (id, power, voltage, date) VALUE (?, ?, ?, ?)",
		id, power, voltage, time.Now().String())
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	return nil
}

func getData(id, start, end string) string {
	var res string
	var rows *sql.Rows
	db, err := sql.Open("mysql",
		SQLConn.Username+":"+SQLConn.Password+"@tcp("+SQLConn.IP+":"+SQLConn.Port+")/plugngo")
	if err != nil {
		log.Println("ERROR:", err)
		return ""
	}
	defer db.Close()
	if len(start) > 0 {
		rows, err = db.Query("SELECT date, power, voltage FROM status WHERE id=? AND date >= ? AND date <= ? ORDER BY date DESC ",
			id, start, end)
		if err != nil {
			log.Println("ERROR:", err)
			return ""
		}
	} else {
		rows, err = db.Query("SELECT date, power, voltage FROM status WHERE id=? AND date >= now() - INTERVAL 10 MINUTE ORDER BY date DESC ",
			id)
		if err != nil {
			log.Println("ERROR:", err)
			return ""
		}
	}
	for rows.Next() {
		var date, power, voltage string
		err = rows.Scan(&date, &power, &voltage)
		//res += date+","+power+","+voltage+"\n"
		res += date+","+power+"\n"
		log.Println(date, power, voltage)
	}
	rows.Close()
	return res
}
