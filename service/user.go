package service

import (
	"../model"
	"../util"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

type UserService struct {
}

//用户注册函数
func (s *UserService) Register(
	mobile,   //手机
	plainpwd, //明文密码
	nickname, //昵称
	avatar, sex string) (user model.User, err error) {
	//检测手机号码是否存在,避免重复注册
	tmp := model.User{}
	_, err = DbEngin.Where("mobile=? ", mobile).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	//如果存在则返回提示已经注册
	if tmp.Id > 0 {
		return tmp, errors.New("该手机号已经注册")
	}
	//否则拼接插入一个新用户数据
	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	//一个随机的salt,配合MD5使用，避免弱口令的存在
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31())
	//用明文和salt通过MD计算出加密信息
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Createat = time.Now()
	//token 可以是一个随机数，用于每次websocket连接鉴权
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())
	//插入一个用户，可能失败
	_, err = DbEngin.InsertOne(&tmp)
	return tmp, err
}

//登录函数
func (s *UserService) Login(
	mobile, //手机
	plainpwd string) (user model.User, err error) {
	//首先通过手机号查询用户
	tmp := model.User{
	}
	DbEngin.Where("mobile = ?", mobile).Get(&tmp)
	if tmp.Id == 0 {
		return tmp, errors.New("该用户不存在")
	}
	if !util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Passwd) {
		return tmp, errors.New("密码不正确")
	}
	//玩家正确登陆，在此刷新一下token,并存到数据库中，用以websocket登陆验证
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)
	tmp.Token = token
	DbEngin.ID(tmp.Id).Cols("token").Update(&tmp)
	return tmp, nil
}

//查找某个用户
func (s *UserService) Find(
	userId int64) (user model.User) {

	//首先通过手机号查询用户
	tmp := model.User{

	}
	DbEngin.ID(userId).Get(&tmp)
	return tmp
}
