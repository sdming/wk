// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package controller

import (
	"errors"
	"fmt"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	"log"
	"strings"
)

func init() {
	userRepo = &UserRepository{
		users: make([]User, 0),
	}

	userRepo.Add(User{
		Name:     "Sheldon",
		Location: "us",
		Gender:   "male",
		Skills:   []string{"Go", "Java"},
	})

	// demo: how to add template func
	wk.TemplateFuncs["location"] = getLocations
	wk.TemplateFuncs["skill"] = getSkills
	wk.TemplateFuncs["locationtext"] = locationText

	boot.Boot(RegisterUserRoute)
}

type User struct {
	Id       int
	Name     string
	Location string
	Gender   string
	Skills   []string
}

type UserRepository struct {
	users []User
}

var userRepo *UserRepository

func (ur *UserRepository) All() []User {
	return ur.users
}

func (ur *UserRepository) Add(u User) int {
	u.Id = len(ur.users) + 1
	ur.users = append(ur.users, u)
	return u.Id
}

func (ur *UserRepository) Update(u User) bool {
	for i := 0; i < len(ur.users); i++ {
		if ur.users[i].Id == u.Id {
			ur.users[i].Gender = u.Gender
			ur.users[i].Location = u.Location
			ur.users[i].Skills = u.Skills
			return true
		}
	}
	return false
}

func (ur *UserRepository) Exists(name string) bool {
	for _, u := range ur.users {
		if strings.EqualFold(u.Name, name) {
			return true
		}
	}
	return false
}

func (ur *UserRepository) Delete(id int) bool {
	for i, u := range ur.users {
		if u.Id == id {
			ur.users = append(ur.users[:i], ur.users[i+1:]...)
			return true
		}
	}
	return false
}

func (ur *UserRepository) Get(id int) (u User, ok bool) {
	for _, u = range ur.users {
		if u.Id == id {
			ok = true
			return
		}
	}
	return
}

func (ur *UserRepository) GetByLocation(location string) []User {
	users := make([]User, 0)
	for _, u := range ur.users {
		if strings.EqualFold(u.Location, location) {
			users = append(users, u)
		}
	}
	return users
}

type ApiRes struct {
	Code    string
	Message string
	Data    interface{}
}

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

//var server *wk.HttpServer
var locations []wk.HtmlOption
var skills []wk.HtmlOption

func RegisterUserRoute(srv *wk.HttpServer) {

	// url: /user/xxx/xxx
	// route to UserController
	// demo how to route use regexp
	srv.RouteTable.Regexp("*", "^/user/?(?P<action>[[:alnum:]]+)?/?(?P<arg>[[:alnum:]]+)?/?").ToController(NewUserController())

	locations = []wk.HtmlOption{}
	srv.Config.AppConfig.MustChild("locations").Value(&locations)

	skills = make([]wk.HtmlOption, 0)
	skillConfig, _ := srv.Config.AppConfig.MustChild("skills").Slice()
	for _, skill := range skillConfig {
		skills = append(skills, wk.HtmlOption{skill, skill})
	}
}

func locationText(v string) string {
	for _, l := range locations {
		if l.Value == v {
			return l.Text
		}
	}
	return v
}

func getLocations() []wk.HtmlOption {
	return locations
}

func getSkills() []wk.HtmlOption {
	return skills
}

// // get: /user/xxx/
// func (uc *UserController) Default(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
// 	return uc.Index(ctx)
// }

// get: /user
func (uc *UserController) Default(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	ctx.ViewData["users"] = userRepo.All()
	return wk.View("user/index.html"), nil
}

// get: /user/location/america
func (uc *UserController) Location(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	location, ok := ctx.RouteData.Str("arg")
	if ok {
		users := userRepo.GetByLocation(location)
		ctx.ViewData["users"] = users
	}
	return wk.View("user/location.html"), nil
}

// get: /user/index/
func (uc *UserController) Index(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	ctx.ViewData["users"] = userRepo.All()
	return wk.View("user/index.html"), nil
}

// get: /user/all/
func (uc *UserController) All(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	ctx.ViewData["users"] = userRepo.All()
	return wk.View("user/_list.html"), nil
}

// get: /user/exists/sheldon
func (uc *UserController) Exists(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	valid := true
	if name := ctx.FV("name"); name != "" {
		valid = !userRepo.Exists(name)
	}
	return wk.Data(valid), nil
}

// get: /user/delete/1
func (uc *UserController) Delete(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	var res ApiRes

	if id, ok := ctx.RouteData.Int("arg"); ok {
		res = ApiRes{
			Code: "",
			Data: userRepo.Delete(id),
		}
	} else {
		res = ApiRes{
			Code:    "1",
			Message: " Invalid Argument",
		}
	}

	return wk.Json(res), nil
}

// get/post: /user/view/1
func (uc *UserController) View(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	var id int
	if id, err = parseId(ctx); err != nil {
		return
	}

	if user, ok := userRepo.Get(id); ok {
		ctx.ViewData["user"] = user
	}
	return wk.View("user/view.html"), nil
}

// get: /user/add/
func (uc *UserController) AddGet(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("user/add.html"), nil
}

// get/post: /user/add
func (uc *UserController) AddPost(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	id := 0
	if name := ctx.FV("name"); name == "" {
		ctx.ViewData["errmsg"] = "name is required"
	} else if name == "error" {
		ctx.ViewData["errmsg"] = "name is error"
	} else if userRepo.Exists(name) {
		ctx.ViewData["errmsg"] = "name exists"
	} else {
		u := User{
			Name:     ctx.FV("name"),
			Location: ctx.FV("location"),
			Gender:   ctx.FV("gender"),
			Skills:   ctx.Request.Form["skill"],
		}
		id = userRepo.Add(u)
	}

	if id <= 0 {
		ctx.ViewData["ctx"] = ctx
		return wk.View("user/add.html"), nil
	}
	return wk.Redirect(fmt.Sprintf("/user/view/%d", id), false), nil
}

func parseId(ctx *wk.HttpContext) (int, error) {
	if id, ok := ctx.RouteData.Int("arg"); !ok {
		return 0, errors.New("Invalid Argument")
	} else {
		return id, nil
	}
}

// get: /user/edit/1
func (uc *UserController) EditGet(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	var id int
	if id, err = parseId(ctx); err != nil {
		return
	}

	if user, ok := userRepo.Get(id); ok {
		ctx.ViewData["user"] = user
	}
	return wk.View("user/edit.html"), nil
}

// post: /user/edit/1
func (uc *UserController) EditPost(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	var id int
	if id, err = parseId(ctx); err != nil {
		return
	}

	u := User{
		Id:       id,
		Location: ctx.FV("location"),
		Gender:   ctx.FV("gender"),
		Skills:   ctx.Request.Form["skill"],
	}
	if userRepo.Update(u) {
		return wk.Redirect(fmt.Sprintf("/user/view/%d", id), false), nil
	}
	return wk.View("user/edit.html"), nil
}

func (uc *UserController) OnActionExecuting(action *wk.ActionContext) {
	log.Println("UserController action executing", action.Context.Request.URL, action.Name)
}

func (uc *UserController) OnActionExecuted(action *wk.ActionContext) {
	log.Println("UserController action executed", action.Context.Request.URL, action.Name, action.Result)
}

func (uc *UserController) OnException(action *wk.ActionContext) {
	log.Println("UserController action executed", action.Context.Request.URL, action.Name, action.Err)
}
