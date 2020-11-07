package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           //消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"`   //谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         //群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     //对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"`     //消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         //预览图片
	Url     string `json:"url,omitempty" form:"url"`         //服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       //简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   //其他和数字相关的
}
//本核心在于形成userid和Node的映射关系
//确保知道每个socket是哪个用户
type Node struct {
	Conn *websocket.Conn
	//存储需要发送的消息
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系表
var clientMap = make(map[int64]*Node, 0)
//读写锁
var rwlocker sync.RWMutex

func Chat(writer http.ResponseWriter,
	request *http.Request) {

	//todo 检验接入是否合法
	query := request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isvalida := checkToken(userId, token)
	//根据token决定要不要把这个连接升级为websocket
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//todo 获得conn,将连接组合成一个数据结构
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//todo 获取用户全部群Id
	comIds := contactService.SearchComunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}
	//todo userid和node形成绑定关系，公共资源并发访问，所以需要上锁
	rwlocker.Lock()
	clientMap[userId] = node
	fmt.Println("now user count is :",len(clientMap))
	rwlocker.Unlock()
	//todo 完成发送逻辑,con
	go sendproc(node)
	//todo 完成接收逻辑
	go recvproc(node)
	//
	sendMsg(userId, []byte("hello,world!"))
}

//发送协程
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

//todo 添加新的群ID到用户的groupset中
func AddGroupId(userId, gid int64) {
	//取得node
	rwlocker.Lock()
	node, ok := clientMap[userId]
	if ok {
		node.GroupSets.Add(gid)
	}
	//clientMap[userId] = node
	rwlocker.Unlock()
	//添加gid到set
}

//接收协程
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		dispatch(data)
		fmt.Printf("recv<=%s", data)
	}
}

//后端调度逻辑处理
func dispatch(data []byte) {
	//todo 解析data为message
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//todo 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		//todo 群聊转发逻辑
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
	case CMD_HEART:
		//todo 一般啥都不做
	}
}

//todo 发送消息
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

//检测是否有效
func checkToken(userId int64, token string) bool {
	//从数据库里面查询并比对
	user := userService.Find(userId)
	return user.Token == token
}
