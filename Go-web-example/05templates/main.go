package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done   bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main(){
	tmpl:=template.Must(template.ParseFiles("layout.html"))
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		data:=TodoPageData{
			PageTitle: "My Todo list",
			Todos:[]Todo{
				{Title:"Task1",Done:false},
				{Title:"Task2",Done:false},
				{Title:"Task3",Done:true},
			},
		}
		tmpl.Execute(w,data)
	})
	http.ListenAndServe(":8080",nil)
}
