package main

import (
	"fmt"
	"strconv"
	"time"
)

import (
	"github.com/AlexStocks/dubbogo/codec/hessian"
	"github.com/AlexStocks/goext/time"
)

type Gender hessian.JavaEnum

const (
	MAN hessian.JavaEnum = iota
	WOMAN
)

var genderName = map[hessian.JavaEnum]string{
	MAN:   "MAN",
	WOMAN: "WOMAN",
}

var genderValue = map[string]hessian.JavaEnum{
	"MAN":   MAN,
	"WOMAN": WOMAN,
}

func (g Gender) JavaClassName() string {
	return "com.ikurento.user.Gender"
}

func (g Gender) String() string {
	s, ok := genderName[hessian.JavaEnum(g)]
	if ok {
		return s
	}

	return strconv.Itoa(int(g))
}

func (g Gender) EnumValue(s string) hessian.JavaEnum {
	v, ok := genderValue[s]
	if ok {
		return v
	}

	return hessian.InvalidJavaEnum
}

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

type HessianUser struct {
	Id   string
	Name string
	Age  int32
	Time time.Time
	Sex  Gender // 注意此处，java enum Object <--> go string
}

func (u HessianUser) String() string {
	return fmt.Sprintf(
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.Id, u.Name, u.Age, u.Time, u.Sex,
	)
}

func (HessianUser) JavaClassName() string {
	return "com.ikurento.user.User"
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

/*
func testCalc() {
	var tcpAddress = "192.168.102.201:20000"

	dubboCtx := &DubboCtx{
		Path:    "com.ikurento.user.UserProvider2", // 注意此处的值是path
		Service: "com.ikurento.user.UserProvider",
		Method:  "Calc",
		Version: "2.0",
		Timeout: 500,
		Return:  reflect.TypeOf(int32(0)),
	}

	var args = []interface{}{int64(1), int64(2)}
	ret, err := SendHession(tcpAddress, dubboCtx, args)
	if err != nil {
		jerrors.ErrorStack(err)
		return
	}
	fmt.Println("ret:", gxlog.ColorSprint(ret))
}

func testSum() {
	var tcpAddress = "127.0.0.1:20000"
	hessian.RegisterPOJO(&Response{})

	dubboCtx := &DubboCtx{
		Path:    "com.ikurento.user.UserProvider2", // 注意此处的值是path
		Service: "com.ikurento.user.UserProvider",
		Method:  "Sum",
		Version: "2.0",
		Timeout: 500,
		Return:  reflect.TypeOf(Response{}),
	}

	var args = []interface{}{int64(1), int64(2)}
	ret, err := SendHession(tcpAddress, dubboCtx, args)
	if err != nil {
		fmt.Printf("Sum(1, 2) = error:%q", jerrors.ErrorStack(err))
		return
	}
	fmt.Println("ret:", gxlog.ColorSprint(ret))
}

*/
