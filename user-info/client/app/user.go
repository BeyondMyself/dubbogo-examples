package main

import (
	"fmt"
)

import (
	"github.com/AlexStocks/goext/time"
)

type JsonRPCUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
	Time int64  `json:"time"`
	Sex  string `json:"sex"`
	// Sex Gender `json:"sex"`
}

func (u JsonRPCUser) String() string {
	return fmt.Sprintf(
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.Id, u.Name, u.Age, gxtime.YMDPrint(int(u.Time), 0), u.Sex,
	)
}

type Response struct {
	Status int
	Err    string
	Data   int
}

func (r Response) String() string {
	return fmt.Sprintf(
		"Response{Status:%d, Err:%s, Data:%d}",
		r.Status, r.Err, r.Data,
	)
}

func (Response) JavaClassName() string {
	return "com.ikurento.user.Response"
}
