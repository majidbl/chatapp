package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  //"fmt"
  "github.com/majidzarephysics/chatapp/model"
  "github.com/gorilla/websocket"
  )
  
  
var clients = make(map[string]*websocket.Conn)
var msgs = make(chan model.Message)
var addr = flag.String("addr", "localhost:8888", "http service address")
    
var upgrader = websocket.Upgrader{} // use default options
    
func echo(w http.ResponseWriter, r *http.Request) {
      clientHandle(w,r)
  }
    
  func clientHandle(w http.ResponseWriter, r *http.Request){
      upgrader.CheckOrigin = func(r *http.Request) bool { return true }
      
      c, err := upgrader.Upgrade(w, r, nil)
      if err != nil {
        log.Print("upgrade:", err)
        return
      }
      defer func() {
        c.Close()
      }()
      for {
        var msg model.Message
        //mt, message, err := c.ReadMessage()
        err := c.ReadJSON(&msg)
        if err != nil {
          log.Println("read:", err)
          break
        }
        senderTemp := msg.Sender
        
        msg.Content = senderTemp + "says -> " + msg.Content
        if msg.Reciver == "server"{
          clients[senderTemp] = c
        }else{
        msgs <- msg
          
        }
      }
  }
func handlemsg(){
  for{
    m := <- msgs
    if _, ok := clients[m.Reciver]; ok{
      err := clients[m.Reciver].WriteJSON(m)
      if err != nil{
        log.Println("write:", err)
        break
      }
    }else{
      fmt.Println("error : message receiver does not existing")
    }
  }
}
    
func hello(w http.ResponseWriter, r *http.Request){
      payload := model.Message{
        Sender:"majid",
        Reciver:"Behnam",
        Content:"test json",
      }
      response, err := json.Marshal(payload)
      if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
      }
      w.Header().Set("Content-Type", "application/json")
      w.Write([]byte(response))
    }
 
 func adduser(w http.ResponseWriter, r *http.Request){
   
      
      var userJson model.Users
      json.NewDecoder(r.Body).Decode(&userJson)
      db, err := model.GetDB()
      if err != nil{
        fmt.Println(err)
        return
      }
      var users []model.Users
      db.Where("user_name = ?", userj.UserName).First(&users)
   
      if len(users) == 0 {
        db.Create(&userj)
      }else{
        fmt.Printf("a user with %s username  Exist", userj.UserName)
      }
 }

func addcontact(w http.ResponseWriter, r *http.Request){
  var userj model.Users
  var userg model.Users
  db, err := model.GetDB()
  if err != nil{
    fmt.Println(err)
    return
  }
  
  
  js := json.NewDecoder(r.Body).Decode(&userj)
  
  db.Where("user_name = ?", userj.UserName).First(&userg)
  db.Find(&userj)
  db.Model(&model.Users{}).Where("user_name = ?", userj.UserName).Update("contact",userg.Contact + "," + userj.Contact)
}

func main() {
      flag.Parse()
      log.SetFlags(0)
      db, err := model.GetDB()
      if err != nil {
        log.Fatal("database initial was failed")
      }
      userTableExist := db.Migrator().HasTable("users")
      if userTableExist{
        fmt.Println("user Table Exist")
      } else {
        db.Migrator().CreateTable(&model.Users{})
      }
      r := mux.NewRouter()
      // Create a simple file server
      //fs := http.FileServer(http.Dir("./frontend"))
      //http.Handle("/", fs)
      r.HandleFunc("/echo", echo)
      r.HandleFunc("/hello", hello)
      r.HandleFunc("/adduser", adduser)
      r.HandleFunc("/addcontact", addcontact)
      go handlemsg()
      http.Handle("/", r)
      log.Fatal(http.ListenAndServe(*addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))

      //log.Fatal(http.ListenAndServe(*addr, nil))
    }
