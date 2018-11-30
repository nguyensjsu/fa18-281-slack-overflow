package main

import (
	"fmt"
	"net/http"
    "log"
    "os"
    "errors"
    // "encoding/json"
    "github.com/gorilla/mux"
	"github.com/gocql/gocql"
	"strconv"
	"github.com/unrolled/render"
	"github.com/codegangsta/negroni"
)

var Session *gocql.Session


type heartbeatResponse struct {
  Status string `json:"status"`
  Code int `json:"code"`
}

type restrauntResponse struct{
	 restraunt_id string `json:"restraunt_id""`  
     name string `json:"name""`
     address string `json:"address""`
     phone string `json:"phone""`
     open_hours string `json:"open_hours""`
}

type menuResponse struct{
	 item_id string `json:"item_id""`  
     name string `json:"name""`
     price float32 `json:"price""`
     restraunt_id string `json:"restraunt_id""`
}

func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

func callFunc() error {
    return errors.New("error")
}

type finalRestrauntResponse struct{
	restrauntResponse[] restrauntResponse
}

type finalMenuResponse struct{
	menuResponse[] menuResponse
}

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	var err error

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "cmpe281"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
	return n
}


func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/restraunts", restrauntHandler(formatter)).Methods("GET")
	mx.HandleFunc("/menus", menuHandler(formatter)).Methods("GET")
	// mx.HandleFunc("/order", gumballNewOrderHandler(formatter)).Methods("POST")
	// mx.HandleFunc("/order/{id}", gumballOrderStatusHandler(formatter)).Methods("GET")
	// mx.HandleFunc("/order", gumballOrderStatusHandler(formatter)).Methods("GET")
	// mx.HandleFunc("/orders", gumballProcessOrdersHandler(formatter)).Methods("POST")
}

func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API is alive!"})
	}
}

func restrauntHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := req.URL.Query()
		var zip int
		zip,_ = strconv.Atoi(vars.Get("zip"))
		var restraunt_name string
		var restraunt_address string
		var phone_number string
		var opening_hours string
		var restraunt_id string
		array := []map[string]interface{}{}
		if err := Session.Query("SELECT restraunt_id,restraunt_name,restraunt_address,phone_number,opening_hours FROM restraunts WHERE zip_code = ? ALLOW FILTERING",zip).Consistency(gocql.One).Scan(&restraunt_id,&restraunt_address,&restraunt_name,&phone_number,&opening_hours); err != nil {
			formatter.JSON(w, http.StatusOK, array)
			return
		}
		iter := Session.Query("SELECT restraunt_id,restraunt_name,restraunt_address,phone_number,opening_hours FROM restraunts WHERE zip_code = ? ALLOW FILTERING",zip).Iter()
		ret := &map[string]interface{}{
			"restraunt_id":     &restraunt_id,
			"restraunt_name": &restraunt_name,
			"restraunt_address": &restraunt_address,
			"phone_number": &phone_number,
			"opening_hours" : &opening_hours,
		}
		for{
			ret = &map[string]interface{}{
				"restraunt_id":     &restraunt_id,
				"restraunt_name": &restraunt_name,
				"restraunt_address": &restraunt_address,
				"phone_number": &phone_number,
				"opening_hours" : &opening_hours,
			}

			if !iter.MapScan(*ret) {
				break
			}
			array = append(array,*ret)
		}
		formatter.JSON(w, http.StatusOK, array)
	}
}


func menuHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := req.URL.Query()
		var res_id string
		res_id = vars.Get("restraunt_id")
		var restraunt_id string
		var item_id string
		var price float32
		var name string
		array := []map[string]interface{}{}
		if err := Session.Query("SELECT item_id,name,price,restraunt_id FROM menu WHERE restraunt_id = ? ALLOW FILTERING",res_id).Consistency(gocql.One).Scan(&item_id,&name,&price,&restraunt_id); err != nil {
			formatter.JSON(w, http.StatusOK, array)
			return
		}
		iter := Session.Query("SELECT item_id,name,price,restraunt_id FROM menu WHERE restraunt_id = ? ALLOW FILTERING",res_id).Iter()

		ret := &map[string]interface{}{
			"restraunt_id":     &restraunt_id,
			"item_id": &item_id,
			"price": &price,
			"name": &name,
		}
		for{
			ret = &map[string]interface{}{
				"restraunt_id":     &restraunt_id,
				"item_id": &item_id,
				"price": &price,
				"name": &name,
			}
			if !iter.MapScan(*ret) {
				break
			}
			array = append(array,*ret)
		}
		formatter.JSON(w, http.StatusOK, array)
	}
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	server := NewServer()
	server.Run(":" + port)
}
