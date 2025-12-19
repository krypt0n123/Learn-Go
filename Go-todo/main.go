package main

import (
	"encoding/json"
	"log"
	"todo/db"

	"github.com/google/uuid"

	"net/http"
)

func main() {

	db.Init()

	//create a new todo
	http.HandleFunc("/create", handleCreateTodo)

	//get all todos
	http.HandleFunc("/get-all-todos", handleGetAllTodos)

	//updata a todo
	http.HandleFunc("/update", handleUpdate)

	//delete a todo
	http.HandleFunc("/delete", handleDelete)

	//处理根路径"/",提供HTML页面
	http.HandleFunc("/", handleIndex)

	//提供静态文件(CSS,JS)
	//1.创建一个文件服务器，指向“dist”目录
	fs := http.FileServer(http.Dir("./frontend/dist"))
	//2.告诉go：任何以“/static/”开头的URL，都从该文件服务器获取
	//	http.StripPrefix会移除“/static/”前缀
	//	所以请求/static/output.css 会变成./frontend/dist 目录查找output.css
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("start server on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	//1.读取前段传来的参数
	params := map[string]string{}
	//2.解析参数
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//3.处理参数
	name := params["name"]
	description := params["description"]
	id := uuid.New().String()
	//4.生成ID
	var newTodo db.Todo = db.Todo{
		ID:          id,
		Name:        name,
		Description: description,
		Completed:   false,
	}
	//5.将数据传入数据库
	err = db.CreateTodo(newTodo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//6.返回结果
	w.WriteHeader(http.StatusOK)
}

func handleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// log.Println("handleGetAllTodos:", db.Todos)
	todos, err := db.GetAllTodos()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := params["id"]
	name := params["name"]
	description := params["description"]
	completed := params["completed"]

	err = db.UpdateTodo(db.Todo{
		ID:          id,
		Name:        name,
		Description: description,
		Completed:   completed == "true",
	})
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusOK)
		return
	}

	id := params["id"]

	err = db.DeleteTodo(id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./frontend/index.html")
}
