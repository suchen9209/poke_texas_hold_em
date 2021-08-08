package models

import "github.com/beego/beego/v2/core/logs"

/**
一整场游戏
**/
type Game struct {
	Id       int    `form:"-"`
	GameName string `form:"name"`
}

/**
一局游戏，从发牌到结算
**/
type GameMatch struct {
	Id                int
	GameId            int
	SmallBindPosition int
	BigBindPosition   int
	ButtonPosition    int
	PotAll            int
	Pot1st            int
	Pot2nd            int
	Pot3rd            int
	Pot4th            int
	GameStatus        string
}

/**
单局游戏玩家操作记录
**/
type GameMatchLog struct {
	Id          int
	GameMatchId int
	UserId      int
	Operation   int
	PointNumber int
}

/**
一整场游戏玩家参与位置，是否在线
**/
type GameUser struct {
	Id       int
	UserId   int
	GameId   int
	Position int
	Online   int
}

const GAMEUSERNUMBER = 8

func SetUserReturnPlayer(u User) int {
	//如果已经存在，则直接返回位置
	var gut = &GameUser{
		UserId: u.Id,
		GameId: 1,
	}
	o.Read(gut, "UserId", "GameId")
	logs.Info(gut)
	if gut.Id > 0 {
		gut.Online = 1
		o.Update(gut)
		return gut.Position
	}

	//如果不存在查找剩余位置，分配座位
	var gu []*GameUser
	number, _ := o.QueryTable("game_user").Filter("game_id", 1).Filter("online", 1).All(&gu)
	if number >= GAMEUSERNUMBER {
		panic("too many player")
	}
	var poslit = make(map[int]int)
	for _, g := range gu {
		poslit[g.Position] = g.UserId

	}
	var pos int
	for i := 1; i <= GAMEUSERNUMBER; i++ {
		if poslit[i] == 0 {
			pos = i
			tmp_g := GameUser{
				UserId:   u.Id,
				GameId:   1,
				Position: pos,
				Online:   1,
			}
			o.Insert(&tmp_g)
			break
		}
	}

	return pos
}

func RemoveGameUser(uid int, gid int) {
	o.QueryTable("game_user").Filter("game_id", gid).Filter("user_id", uid).Delete()
}
