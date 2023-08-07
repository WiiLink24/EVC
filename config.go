package main

import (
	"encoding/xml"
	"io/ioutil"
)

type Config struct {
	Username        string `xml:"username"`
	Password        string `xml:"password"`
	DatabaseAddress string `xml:"databaseAddress"`
	DatabaseName    string `xml:"databaseName"`
}

func GetConfig() Config {
	data, err := ioutil.ReadFile("config.xml")
	checkError(err)

	var config Config
	err = xml.Unmarshal(data, &config)
	checkError(err)

	return config
}
