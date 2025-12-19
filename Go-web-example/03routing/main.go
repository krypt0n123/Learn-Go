package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//处理网站根目录
func HomeHandler(w http.ResponseWriter, r*http.Request){
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w,"welcom to E-commerce API homepage")
}

//处理获取单个商品的详情. /products/{id}
func ProductDetailHandler(w http.ResponseWriter ,r *http.Request){
	//使用mux.Vars(r)来获取一个map，其中包含了URL中的所有变量
	vars:=mux.Vars(r)
	id:=vars["id"]

	fmt.Fprintf(w,"You are viewing details for Product ID:%s\n",id)
}

//多个路径变量：处理获取某个用户的订单
func UserOrderDetailHandler(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	username:=vars["username"]
	orderID:=vars["order_id"]

	fmt.Fprintf(w,"User '%s' is viewing Order ID:%s\n",username,orderID)
}

//限制方法（GET）：来获取所有商品列表
func ProductListHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"Listing all products(GET)")
}

//限制方法（POST）：创建一个新商品
func ProductCreateHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"Creating a new product(POST)")
}

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"welcom to the ADMIN Dashboard!")
	fmt.Fprintln(w,"You must be accessing this form 'admin.example.com'")
}

func PaymentHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"Processing secure payment(HTTPS only)")
}

func main(){
	//创建一个新的路由实例
	r:=mux.NewRouter()

	r.HandleFunc("/",HomeHandler).Methods("GET")
	//所有API路由都在/api/v1/前缀下，方便对API进行版本控制
	api:=r.PathPrefix("/api/v1").Subrouter()

	//GET /api/v1/products - 获取列表
	api.HandleFunc("/products",ProductListHandler).Methods("GET")
	//POST /api/v1/products - 创建列表
	api.HandleFunc("/products",ProductCreateHandler).Methods("POST")

	//GET /api/v1/products/{id} - 获取单个商品
	api.HandleFunc("/products/{id:[0-9]+}",ProductDetailHandler).Methods("GET")

	//GET /api/v1/users/{username}/orders/{order_id}
	//{username}可以是字母数字, {order_id}必须是数字
	api.HandleFunc("/users/{username:[a-zA-Z0-9]+}/orders/{order_id:[0-9]+}",UserOrderDetailHandler).Methods("GET")

	//POST /api/v1/payment
	//这个路由只接受HTTPS请求，用于处理敏感的支付信息
	api.HandleFunc("/payment",PaymentHandler).Methods("POST").Schemes("https")

	//GET /dashboard
	//只在Host头部为“admin/example.com”时才会匹配
	r.HandleFunc("/dashboard",AdminDashboardHandler).Host("admin.example.com")

	log.Println("Server start on localhost:8080")

	err:=http.ListenAndServe(":8080",r)
	if err!=nil{
		log.Fatal(err)
	}
}