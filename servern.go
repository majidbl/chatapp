package main

import (
  "flag"
  "fmt"
  "log"
  "strings"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "reflect"
  "github.com/majidzarephysics/chatapp/model"
  "github.com/gorilla/websocket"
  )
  
  
var clients = make(map[string]*websocket.Conn)
var msgs = make(chan model.Message)
var addr = flag.String("addr", "localhost:8888", "http service address")
type MsgArr struct{
  messagesArray []model.Message
}   
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
        
        //msg.Content = senderTemp + "says -> " + msg.Content
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
    errm := model.SaveMsgDB(m)
    if errm != nil{
      fmt.Println(errm)
    }
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
 func getMsg(w http.ResponseWriter, r *http.Request){
   sr := mux.Vars(r)
   s_r := strings.Split(sr["sender_reciver"],"_")
   db, errdb := model.GetDB()
   if errdb != nil{
     fmt.Println(errdb)
     return
   }
   //fmt.Println(sr[0])
   //fmt.Println(sr[1])
   var msgarr []model.Message
   db.Where("sender=?",s_r[0]).Where("reciver= ?",s_r[1]).Or(db.Where("sender=?",s_r[1]).Where("reciver= ?",s_r[0])).Find(&msgarr)
   response, err := json.Marshal(msgarr)
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
      db.Where("user_name = ?", userJson.UserName).First(&users)
   
        fmt.Println(reflect.TypeOf(userJson))
      if len(users) == 0 {
        db.Create(&userJson)
      }else{
        fmt.Printf("a user with %s username  Exist", userJson.UserName)
      }
 }

func addcontact(w http.ResponseWriter, r *http.Request){
  var userJson model.Users
  var userg model.Users
  json.NewDecoder(r.Body).Decode(&userJson)
  db, err := model.GetDB()
  
  if err != nil{
    fmt.Println(err)
    return
  }
  
  db.Where("user_name = ?", userJson.UserName).First(&userg)
  fmt.Println("dats from frontend",userJson)
  fmt.Println("dats from backend",userg)
  db.Model(&model.Users{}).Where("user_name = ?", userJson.UserName).Update("contact",userg.Contact + userJson.Contact + ",")
  
  
}
func getContact(w http.ResponseWriter, r *http.Request){
  un := mux.Vars(r)
  var userc model.Users
  db,err := model.GetDB()
  if err != nil{
    fmt.Println(err)
  }
  db.Where("user_name = ?", un["user_name"]).First(&userc)
  fmt.Println(un["user_name"])
  response, err := json.Marshal(userc)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
      }
      w.Header().Set("Content-Type", "application/json")
      w.Write([]byte(response))
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
      msgTableExist := db.Migrator().HasTable("messages")
      if msgTableExist{
        fmt.Println("Message Table Exist")
      } else {
        db.Migrator().CreateTable(&model.Message{})
      }
      r := mux.NewRouter()
      // Create a simple file server
      //fs := http.FileServer(http.Dir("./frontend"))
      //http.Handle("/", fs)
      r.HandleFunc("/echo", echo)
      r.HandleFunc("/hello", hello)
      r.HandleFunc("/adduser", adduser)
      r.HandleFunc("/addcontact", addcontact)
      r.HandleFunc("/getcontact/{user_name}", getContact)
      r.HandleFunc("/getmessage/{sender_reciver}", getMsg)
      go handlemsg()
      http.Handle("/", r)
      log.Fatal(http.ListenAndServe(*addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "XMLHttpRequest", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))

      //log.Fatal(http.ListenAndServe(*addr, nil))
    }
