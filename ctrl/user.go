package ctrl

import (
	"../model"
	"../service"
	"../util"
	"fmt"
	"math/rand"
	"net/http"
)
var userService service.UserService
func UserLogin(writer http.ResponseWriter,
	request *http.Request) {
	request.ParseForm()
	//获取手机号和密码信息
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	//校验需要用到数据库，交给service层完成
	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}

}
func UserRegister(writer http.ResponseWriter,
	request *http.Request) {
	request.ParseForm()
	//注册也需要手机号和密码
	mobile := request.PostForm.Get("mobile")
	//明文密码，存储到数据库时需要md5加密
	plainpwd := request.PostForm.Get("passwd")
	//随机给一个用户名吧,正常应该不重复且按照一定次序来
	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := model.SEX_UNKNOW
	//这个也需要用到数据库操作故交给service层操作
	user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}
}
