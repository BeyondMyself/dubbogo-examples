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
}

func (u JsonRPCUser) String() string {
	return fmt.Sprintf(
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.Id, u.Name, u.Age, gxtime.YMDPrint(int(u.Time), 0), u.Sex,
	)
}
