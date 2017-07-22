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
