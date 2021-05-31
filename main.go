package main

import (
	"NatsToKafka/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAllChannels() ([]models.Channel,error) {
	client := &http.Client{}
	var channels []models.Channel

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8888/api-gw/v1/thing/getall", nil)
	if err !=nil{
		return channels,err
	}
	req.Header.Set("Authorization","Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJOb2RlUmVkVXNlckBtb2JpZm9uZS52biIsImlhdCI6MTYyMjI2MTk2OSwiZXhwIjoxNjIyMzQ4MzY5fQ.2GDbnW8A-U1DKN_LB6l69tt1aXQfBN49JxMwoX9Jgg2K10Kd0AnQuPRU1i9voZ-D4X87ZNYPh0oTmeRZAIkCYQ")
	req.Header.Set("Accept","application/json")
	resTTTM, err := client.Do(req)
	if err!=nil{
		return channels,err
	}
	data, _ := ioutil.ReadAll(resTTTM.Body)
	err = json.Unmarshal(data,&channels)
	if err!=nil{
		return channels,err
	}
	return channels,nil
}
func main() {
	channels,err:=GetAllChannels()
	if err==nil{
		fmt.Println(channels[0].Id)
	}
}