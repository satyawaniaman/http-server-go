package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)
type User struct{
	Name string `json:"name"`
}
var userCache = make(map[int]User)

var cacheMutex sync.RWMutex
func main(){
	mux:=http.NewServeMux() // create a new serve mux and handles requests
	fmt.Println("Server is running on port 8080")
	mux.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("Hello World"))
	})
	mux.HandleFunc("POST /users",func(w http.ResponseWriter,r *http.Request){
		var user User
		err:= json.NewDecoder(r.Body).Decode(&user)
		if err!= nil{
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}
		if user.Name == ""{
			http.Error(w,"Name is required",http.StatusBadRequest)
			return
		}
		cacheMutex.Lock()
		userCache[len(userCache)+1]=user
		cacheMutex.Unlock()
		w.Write([]byte("User Created"))
	})
	http.ListenAndServe(":8080",mux)
	mux.HandleFunc("GET /users/{id}",func(w http.ResponseWriter,r *http.Request){
		id,err:= strconv.Atoi(r.PathValue("id"))
		if err!= nil{
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}
		cacheMutex.RLock()
		user,ok:=userCache[id]
		cacheMutex.RUnlock()
		if !ok{
			http.Error(w,"User not found",http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type","application/json")
		j,err:= json.Marshal(user)
		if err!= nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	})
	mux.HandleFunc("DELETE /users/{id}",func(w http.ResponseWriter,r *http.Request){
		id,err:= strconv.Atoi(r.PathValue("id"))
		if err!= nil{
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}
		if _,ok:=userCache[id];!ok{
			http.Error(w,"User not found",http.StatusNotFound)
			return
		}
		cacheMutex.Lock()
		delete(userCache,id)
		cacheMutex.Unlock()
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("User Deleted"))
	})
}