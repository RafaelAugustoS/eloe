package main

import (
      "bytes"
      "time"
      "net/http"
      "encoding/json"
      "github.com/aichaos/rivescript-go"
)

var eloe *rivescript.RiveScript = rivescript.New(rivescript.WithUTF8())

var client *http.Client = &http.Client{}

var userId string

var token string = ""
var matchId string = ""

func getUserId() {
      type User struct {
            Id string `json:"_id"`
      }
      var user User
      req, err := http.NewRequest("GET", "https://api.gotinder.com/profile", nil)
      if err != nil {
            panic(err)
      }
      req.Header.Set("x-auth-token", token)
      res, err := client.Do(req)
      if err != nil {
            panic(err)
      }
      json.NewDecoder(res.Body).Decode(&user)
      userId = user.Id
}

func getMessage() {
      type Message struct {
            Data struct {
                  Messages []struct {
                        Message string `json:"message"`
                        To string `json:"to"`
                  } `json:"messages"`
            } `"json:"data"`
      }
      var message Message
      req, err := http.NewRequest("GET", "https://api.gotinder.com/v2/matches/" + matchId + "/messages?count=1", nil)
      if err != nil {
            panic(err)
      }
      req.Header.Set("x-auth-token", token)
      res, err := client.Do(req)
      if err != nil {
            panic(err)
      }
      json.NewDecoder(res.Body).Decode(&message)
      if len(message.Data.Messages) == 0 {
            talk("start a conversation")
      } else {
            if message.Data.Messages[0].To == userId {
                  talk(message.Data.Messages[0].Message)
            }
      }
}

func main() {
      getUserId()
      err := eloe.LoadDirectory("brain")
      if err != nil {
            panic(err)
      }
      eloe.SortReplies()
      ticker := time.NewTicker(time.Second * 10)
      for range ticker.C {
            getMessage()
      }
}

func talk(message string) {
      reply, err := eloe.Reply("local-user", message)
      if err != nil {
            panic(err)
      } else {
            postMessage(reply)
      }
}

func postMessage(reply string) {
      type Request struct {
		Message string `json:"message"`
	}
      var request Request
      request.Message = reply
      body, err := json.Marshal(request)
      if err != nil {
            panic(err)
      }
      req, err := http.NewRequest("POST", "https://api.gotinder.com/user/matches/" + matchId, bytes.NewBuffer(body))
      if err != nil {
            panic(err)
      }
      req.Header.Set("x-auth-token", token)
      req.Header.Set("Content-Type", "application/json")
      _, err = client.Do(req)
      if err != nil {
            panic(err)
      }
}
