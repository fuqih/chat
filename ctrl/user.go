package ctrl

import (
	"../model"
	"../service"
	"../util"
	"net/http"
)

func UserLogin(writer http.ResponseWriter,
	request *http.Request) {
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}

}

var userService service.UserService

func UserRegister(writer http.ResponseWriter,
	request *http.Request) {
	//数据绑定
	var user model.User
	util.Bind(request, &user)
	user, err := userService.Register(user.Mobile, user.Passwd, user.Nickname, user.Avatar, user.Sex)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}
}
