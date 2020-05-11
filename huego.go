package main

import (
  "os"
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "strings"
  "flag"
)

var USER string // get user string here: https://developers.meethue.com/develop/get-started-2/
var IPADDRESS string // the IP address of your hue bridge

// add your group ID mappings here
var All int = 0
var DiningRoom int = 1
var Bedroom int = 3
var LivingRoom int = 4

type Group struct {
  Name    string    `json:"name"`
  State   State     `json:"state"`
}

type State struct {
  AllOn    bool    `json:"all_on"`
  AnyOn    bool    `json:"any_on"`
}

func main() {
  USER = os.Getenv("USER")
  IPADDRESS = os.Getenv("IPADDRESS")

  var flagGroup = flag.String("group", "all", "Group to use: all, diningroom, bedroom, livingroom")
  var flagToggle = flag.Bool("toggle", false, "Toggle group: true / false")
  var flagHelp = flag.Bool("help", false, "Show usage")
  flag.Parse()

  if (*flagToggle == true){
    // add your groups here
    if (*flagGroup == "all"){ toggleGroup(All) }
    if (*flagGroup == "diningroom"){ toggleGroup(DiningRoom) }
    if (*flagGroup == "bedroom"){ toggleGroup(Bedroom) }
    if (*flagGroup == "livingroom"){ toggleGroup(LivingRoom) }
  }

  if (*flagHelp == true){
    fmt.Println("HueGo usage:")
    fmt.Println("USER='[your hue user id]' IPADDRESS='[hue bridge ip address]' huego --toggle --group=bedroom")
  }
}

func toggleGroup(group int){
  var anyOn bool = getGroupState(group)
  setGroupState(group, !anyOn)
}

func getGroupState(group int) bool {
  var url string = "http://" + IPADDRESS + "/api/" + USER + "/groups/" + strconv.Itoa(group)
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  var currentGroup Group
  json.Unmarshal(body, &currentGroup)
  var anyOn bool = currentGroup.State.AnyOn;

  return anyOn
}

func setGroupState(group int, state bool) {
  var url string = "http://" + IPADDRESS + "/api/" + USER + "/groups/" + strconv.Itoa(group) + "/action/"

  client := &http.Client{}

  body := fmt.Sprintf("{\"on\":%t}", state)

  req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
  if err != nil {
    panic(err)
  }

  req.Header.Set("Content-Type", "application/json; charset=utf-8")
  _, err = client.Do(req)
  if err != nil {
    panic(err)
  }
}
