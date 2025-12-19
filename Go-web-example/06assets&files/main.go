package main

import (
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // <-- 正确的

	http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		http.ServeFile(w,r,"assets/index.html")
	})
	
	http.ListenAndServe(":8080", nil)
}
