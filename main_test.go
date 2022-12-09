package main

import (
	"fmt"
	"testing"
)

type DataList struct {
	List []Item `yaml:"list"`
}

type Item struct {
	ProjectName string `yaml:"project-name"`
	Number      string `yaml:"number"`
}

func TestResult(*testing.T) {
	var dataList DataList
	dataList.List = make([]Item, 0)
	dataList.List = append(dataList.List, Item{
		ProjectName: "admin",
		Number:      "1",
	})
	dataList.List = append(dataList.List, Item{
		ProjectName: "root",
		Number:      "2",
	})
	dataList.List = append(dataList.List, Item{
		ProjectName: "root",
		Number:      "3",
	})

	for i, v := range dataList.List {
		fmt.Println("before: ", v)
		if i == 1 {
			v.Number = "upgrade"
		}
		fmt.Println("after: ", v)
	}
}

func TestBreak(*testing.T) {
	data := make([]string, 4)
	data[0] = "admin"
	data[1] = "root"
	data[2] = "guest"
	data[3] = "test"

	for i, v := range data {
		fmt.Println(i, v)
		if i == 2 {
			break
		}
	}
}
