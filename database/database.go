package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Proxy struct {
	ID           int
	Protocol     string
	SrcAddr      string
	DestAddr     string
	ResponseBody string
	ResponseCode string
	HexStr       string
}

func Init() *sql.DB {
	// 创建数据库连接
	db, err := sql.Open("sqlite3", "./proxys.db")
	if err != nil {
		log.Fatal(err)
	}

	// 创建表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS proxys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		protocol TEXT,
		srcaddr TEXT NOT NULL,
		destaddr TEXT,
		responsebody TEXT,
		responsecode TEXT,
		hexstr TEXT,
		UNIQUE (srcaddr)
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InsertDB(db *sql.DB, protocol, srcaddr, destaddr, responsebody, responsecode, hexstr *string) {
	// 插入数据
	stmt, err := db.Prepare("INSERT INTO proxys(protocol, srcaddr, destaddr, responsebody, responsecode, hexstr) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(*protocol, *srcaddr, *destaddr, *responsebody, *responsecode, *hexstr)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Insert the %d data\n", id)
}

func QuertDB(db *sql.DB, srcaddr *string) *[]Proxy {
	// 查询数据
	var sql string
	if srcaddr != nil {
		sql = fmt.Sprintf("SELECT * FROM proxys WHERE srcaddr = \"%s\"", *srcaddr)
	} else {
		sql = "SELECT * FROM proxys"
	}
	// fmt.Println("test_sql", sql)
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var proxys []Proxy
	for rows.Next() {
		var proxy Proxy
		err = rows.Scan(&proxy.ID, &proxy.Protocol, &proxy.SrcAddr, &proxy.DestAddr, &proxy.ResponseBody, &proxy.ResponseCode, &proxy.HexStr)
		if err != nil {
			log.Fatal(err)
		}
		proxys = append(proxys, proxy)
	}
	// log.Println("test_eeeee: proxys=", proxys)
	return &proxys
}

func QuertNumDB(db *sql.DB, hexstr string) *[]Proxy {
	// 查询指定id数据
	sql := fmt.Sprintf("SELECT * FROM proxys WHERE hexstr = \"%s\"", hexstr)

	// fmt.Println("test_sql", sql)
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var proxys []Proxy
	for rows.Next() {
		var proxy Proxy
		err = rows.Scan(&proxy.ID, &proxy.Protocol, &proxy.SrcAddr, &proxy.DestAddr, &proxy.ResponseBody, &proxy.ResponseCode, &proxy.HexStr)
		if err != nil {
			log.Fatal(err)
		}
		proxys = append(proxys, proxy)
	}
	// log.Println("test_eeeee: proxys=", proxys)
	return &proxys
}

func UpdataDB(db *sql.DB, protocol, destaddr, responsebody, responsecode, srcaddr, hexstr *string) {
	// 更新数据
	stmt, err := db.Prepare("UPDATE proxys SET protocol=?, destaddr=?, responsebody=?, responsecode=?, hexstr=? WHERE srcaddr=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(*protocol, *destaddr, *responsebody, *responsecode, *hexstr, *srcaddr)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Update the %d date\n", rowsAffected)
}

func DeleteDB(db *sql.DB, num int) {
	// 删除数据
	stmt, err := db.Prepare("DELETE FROM proxys WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(num)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Delete the %d data\n", rowsAffected)
}

func DeleteDbHexStr(db *sql.DB, hexstr string) {
	// 删除数据
	stmt, err := db.Prepare("DELETE FROM proxys WHERE hexstr=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(hexstr)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Delete the %d data\n", rowsAffected)
}
