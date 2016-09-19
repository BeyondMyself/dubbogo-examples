/*****************************************************
# DESC    : UserProvider Service
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : LGPL V3
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-07-21 19:22
# FILE    : user.go
******************************************************/

package main

import (
	// "encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

import (
	"github.com/AlexStocks/gocolor"
)

import (
//"github.com/AlexStocks/dubbogo/common"
)

type Gender int

const (
	MAN = iota
	WOMAN
)

var genderStrings = [...]string{
	"man",
	"woman",
}

func (g Gender) String() string {
	return genderStrings[g]
}

type (
	User struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Age   int    `json:"age"`
		sex   Gender
		Sex   string `json:"sex"`
		Birth int    `json:"time"`
	}

	UserId struct {
		Id string
	}

	UserProvider struct {
		user map[string]User
	}
)

var (
	DefaultUser = User{
		Id: "0", Name: "Alex Stocks", Age: 31,
		Birth: int(time.Date(1985, time.November, 10, 23, 0, 0, 0, time.UTC).Unix()),
		sex:   Gender(MAN),
	}

	userMap = UserProvider{user: make(map[string]User)}
)

func init() {
	DefaultUser.Sex = DefaultUser.sex.String()
	userMap.user["A001"] = User{Id: "001", Name: "ZhangSheng", Age: 18, sex: MAN}
	userMap.user["A002"] = User{Id: "002", Name: "Lily", Age: 20, sex: WOMAN}
	userMap.user["A003"] = User{Id: "113", Name: "Moorse", Age: 30, sex: MAN}
	for k, v := range userMap.user {
		v.Birth = int(time.Now().AddDate(-1*v.Age, 0, 0).Unix())
		v.Sex = userMap.user[k].sex.String()
		userMap.user[k] = v
	}
}

/*
// you can define your json unmarshal function here
func (this *UserId) UnmarshalJSON(value []byte) error {
	this.Id = string(value)
	this.Id = common.TrimPrefix(this.Id, "\"")
	this.Id = common.TrimSuffix(this.Id, `"`)

	return nil
}
*/

func (this *UserProvider) getUser(userId string) (*User, error) {
	if user, ok := userMap.user[userId]; ok {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid user id:%s", userId)
}

/*
// can not work
func (this *UserProvider) GetUser(ctx context.Context, req *UserId, rsp *User) error {
	var (
		err  error
		user *User
	)
	user, err = this.getUser(req.Id)
	if err == nil {
		*rsp = *user
		gocolor.Info("rsp:%#v", rsp)
		// s, _ := json.Marshal(rsp)
		// fmt.Println(string(s))

		// s, _ = json.Marshal(*rsp)
		// fmt.Println(string(s))
	}
	return err
}
*/

/*
// work
func (this *UserProvider) GetUser(ctx context.Context, req *string, rsp *User) error {
	var (
		err  error
		user *User
	)

	gocolor.Info("req:%#v", *req)
	user, err = this.getUser(*req)
	if err == nil {
		*rsp = *user
		gocolor.Info("rsp:%#v", rsp)
		// s, _ := json.Marshal(rsp)
		// fmt.Println(string(s))

		// s, _ = json.Marshal(*rsp)
		// fmt.Println(string(s))
	}
	return err
}
*/

func (this *UserProvider) GetUser(ctx context.Context, req []string, rsp *User) error {
	var (
		err  error
		user *User
	)

	gocolor.Info("req:%#v", req)
	user, err = this.getUser(req[0])
	if err == nil {
		*rsp = *user
		gocolor.Info("rsp:%#v", rsp)
		// s, _ := json.Marshal(rsp)
		// fmt.Println(string(s))

		// s, _ = json.Marshal(*rsp)
		// fmt.Println(string(s))
	}
	return err
}

func (this *UserProvider) Service() string {
	return "com.youni.UserProvider"
}

func (this *UserProvider) Version() string {
	return ""
}
