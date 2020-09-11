package model

import "github.com/jinzhu/gorm"

type Users struct{
  gorm.Model
  UserName string `json:"username"`
  Contact  string `json: "contact"`
}

type Message struct {
  gorm.Model
  Sender  string `json:"sender"`
  Reciver string `json:"reciver"`
  Content string `json:"content"`
  }
func (msg *Message) TableName(t string) string {
	return "message"
}
