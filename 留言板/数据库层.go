package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/dome7?charset=utf8")
	if err != nil {
		fmt.Println("sql.open is error !", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	err = db.Ping()
	if err != nil {
		fmt.Println("db.ping is error", err)
	}
}

func Conn() *sql.DB {
	return db
}

func Close() interface{} {
	return db.Close()
}

func Loginll(username string, password string) bool {
	db := Conn()
	res, err := db.Query("select id from user where username=? and password=? ", username, password)
	if err != nil {
		fmt.Println("db.query is error login:", err)
	}

	for res.Next() {
		var id int
		err := res.Scan(&id)
		fmt.Println("id =======",id)
		if err != nil {
			fmt.Println("res.Scan is error:", err)
		}

		if id >= 0 {
			return true
		} else {
			return false
		}
	}
	return false
}

func UserSignup(username string, passwd string,id string) bool {//注册是否成功
	stmt,err := db.Prepare("insert into user(username,password,id) values (?,?,?)")
	if err != nil{
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(username, passwd,id)
	if err != nil {
		fmt.Println("Failed to insert, err:" ,err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func UserExist(username string) bool {//用户是否存在
	var passwd string
	err := db.QueryRow(`select password from user where username = ?`,username).Scan(&passwd)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	if passwd==""{return true}
	return false
}

func InformationSignup(username string,gender string,introduce string ,address string,industry string,occupation string,education string)bool{
	stmt,err := db.Prepare("insert into information(username,gender,introduce,address,industry,occupation,education) values (?,?,?,?,?,?,?)")
	if err != nil{
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(username,gender,introduce,address,industry,occupation,education)
	if err != nil {
		fmt.Println("Failed to insert, err:" ,err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func InformationSelect(username string)(gender,introduce,address,industry,occupation,education string){
	err := db.QueryRow(`select gender,introduce,address,industry,occupation,education from information where username = ?`,username).Scan(&gender,&introduce,&address,&industry,&occupation,&education)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return gender,introduce,address,industry,occupation,education
}

func FindMessageByPid(pid int) []Message {
	rows, err := db.Query("select id,message,user_id from message where pid=?", pid)
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var messageSlice []Message
	for rows.Next() {
		var messages Message
		err := rows.Scan(&id, &messages.Message, &messages.UserId)
		if  err != nil {
			log.Fatal(err)
		}
		child := FindMessageByPid(id)
		messages.ChildMessage = &child
		messageSlice = append(messageSlice, messages)
		//fmt.Println("数据库装载数据时我们的messages", messages, "----messageSlice:", messageSlice)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return messageSlice
}

func MessageInsert(pid int,id int,user_id int,message string) bool {//注册是否成功
	stmt,err := db.Prepare("insert into message(pid,id,user_id,message) values (?,?,?,?)")
	if err != nil{
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(pid,id,user_id,message)
	if err != nil {
		fmt.Println("Failed to insert, err:" ,err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func SelectIdMax()int{
	var id int
	err := db.QueryRow("select MAX(id) from message").Scan(&id)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return id
}

func DeleteMesage(id int)bool{
	res, err := db.Exec("delete from message where id = ? ",id)
	fmt.Println(id)
	if err != nil{
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}