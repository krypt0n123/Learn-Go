package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
}

func main() {
	// 连接数据库
	// 根据 Docker 配置生成的数据源名称 (DSN)
	// 用户: root
	// 密码: my-secret-pw
	// 地址: 127.0.0.1:3306
	// 数据库: mygodo
	// 参数: parseTime=true
	dsn := "root:my-secret-pw@(127.0.0.1:3306)/mygodo?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("can't connect database: ", err)
	}
	defer db.Close() //确保程序退出时关闭连接

	//验证连接是否真正建立
	if err := db.Ping(); err != nil {
		log.Fatal("can't connect database(ping lose): ", err)
	}
	fmt.Println("success connect 'mygodo'database in Docker")

	//创建表
	{
		query := `
		CREATE TABLE IF NOT EXISTS users(
			id INT AUTO_INCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME,
			PRIMARY KEY(id)
		);`
		if _, err := db.Exec(query); err != nil {
			log.Fatal("create table fail", err)
		}
		fmt.Println("'users'table is already")
	}

	//清理表
	{
		if _, err := db.Exec(`DELETE FROM users`); err != nil {
			log.Fatal("clean table fail", err)
		}
		fmt.Println("'users'table is empty")
	}

	//插入一个新用户
	var insertedID int64
	{
		username := "johndoe"
		password := "secret"
		createdAt := time.Now()

		query := `INSERT INTO users (username,password,created_at) VALUES (?,?,?)`
		result, err := db.Exec(query, username, password, createdAt)
		if err != nil {
			log.Fatal("insert 'johndoe' fail:", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatal("get LastInsertId fail:", err)
		}
		insertedID = id
		fmt.Printf("success insert 'johdoe', ID: %d\n", insertedID)
	}

	//插入第二个用户
	{
		query := `INSERT INTO users (username,password,created_at) VALUES (?,?,?)`
		if _, err := db.Exec(query, "alice", "supersecret", time.Now()); err != nil {
			log.Fatal("insert 'alice' fail:", err)
		}
		fmt.Println("success insert second user")
	}

	//查询单个用户
	{
		var (
			id        int
			username  string
			password  string
			createdAt time.Time
		)

		query := "SELECT id, username,password,created_at FROM users WHERE id = ?"
		//使用我们之前保存的insertedID
		err := db.QueryRow(query, insertedID).Scan(&id, &username, &password, &createdAt)
		if err != nil {
			log.Fatal("select fail:", err)
		}
		fmt.Printf("成功查询到单个用户:\n  ID: %d\n  Username: %s\n  Password: %s\n  Created: %v\n", id, username, password, createdAt)
	}

	//查询所有用户
	{
		rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
		if err != nil {
			log.Fatal("select all useers fail: ", err)
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt); err != nil {
				log.Fatal("scan user row fail:", err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal("Row 迭代出错:", err)
		}
		fmt.Println("success select all user:")
		fmt.Printf("%#v\n,users")
	}
	{
		query := `DELETE FROM users WHERE id = ?`
		result, err := db.Exec(query, insertedID) //delete johndoe
		if err != nil {
			log.Fatal("delete user fail:", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatal("get RowsAffected fail: ", err)
		}
		fmt.Printf("success delete user! %d行受到影响 \n", rowsAffected)
	}
	fmt.Println("🎉 练习完成!")
}
