package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Welcome Page</title>

			<link rel="stylesheet" href="/static/css/style.css">

		</head>
		<body>

			<div class="container">
				<h1>Welcome to my website!</h1>
			</div>

		</body>
		</html>
		`)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server statring on http://localhost:8080")

	err:=http.ListenAndServe(":8080", nil)
	if err!=nil{
		log.Fatal("ListenAndServer:",err)
		return
	}
}
