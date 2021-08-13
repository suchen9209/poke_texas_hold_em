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
	Operation   string
	PointNumber int
}

type UserOperationMsg struct {
	Type         EventType
	GameMatchLog GameMatchLog
	Position     int
	Name         string
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

type UserPointSeat struct {
	UserId   int
	Name     string
	Position int
	Point    int
}

const GAME_OP_RAISE = "raise"
const GAME_OP_CALL = "call"
const GAME_OP_CHECK = "check"
const GAME_OP_FOLD = "fold"
const GAME_OP_ALLIN = "allin"
const GAMEUSERNUMBER = 8

func SetUserReturnPlayer(u User) GameUser {
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
		return *gut
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
	var tmp_g GameUser
	var pos int
	for i := 1; i <= GAMEUSERNUMBER; i++ {
		if poslit[i] == 0 {
			pos = i
			tmp_g.UserId = u.Id
			tmp_g.GameId = 1
			tmp_g.Position = pos
			tmp_g.Online = 1
			o.Insert(&tmp_g)
			break
		}
	}

	return tmp_g
}

func RemoveGameUser(uid int, gid int) {
	o.QueryTable("game_user").Filter("game_id", gid).Filter("user_id", uid).Delete()
}

func GetUserPointWithSeat(gid int) []UserPointSeat {
	var gu []UserPointSeat
	o.Raw(`select gu.user_id,u.name,gu.position, u.point from game_user as gu left join user as u on u.id = gu.user_id where gu.game_id = ?`).SetArgs(1).QueryRows(&gu)
	return gu
}

func InitGameMatch(gid int) GameMatch {
	var tmp, new_game GameMatch
	var game_user_list []GameUser
	o.QueryTable("game_user").Filter("game_id", gid).Filter("online", 1).All(&game_user_list)
	if game_user_list == nil {
		panic("no user")
	}
	tmpMap := make(map[int]int)
	for _, v := range game_user_list {
		tmpMap[v.Position] = v.UserId
	}
	new_game.GameId = gid
	o.QueryTable("game_match").Filter("game_id", gid).OrderBy("-id").One(&tmp)
	if tmp.Id <= 0 {
		//之前无对局
		new_game.SmallBindPosition = getPosition(0, tmpMap)

	} else {
		new_game.SmallBindPosition = getPosition(tmp.SmallBindPosition, tmpMap)
	}
	new_game.BigBindPosition = getPosition(new_game.SmallBindPosition, tmpMap)
	new_game.GameStatus = "INIT"
	o.Insert(&new_game)
	return new_game

}

func UpdateGameMatchStatus(gm GameMatch, key string) {
	o.Update(&gm, key)
}

func getPosition(p1 int, tmap map[int]int) int {
	for i := p1 + 1; i <= p1+8; i++ {
		key := i % 8
		if key == 0 {
			key = 8
		}
		if _, ok := tmap[key]; ok {
			return key
		}
	}
	return 0
}

type InRoundUserDetail struct {
	UserId     int
	Position   int
	Point      int
	Name       string
	RoundPoint int      //本轮押注
	AllowOp    []string //raise
}

func GetRoundUserDetail(gid int) []InRoundUserDetail {
	var rmap []InRoundUserDetail
	o.Raw("SELECT gu.`user_id`,gu.`position`,u.point,u.name FROM game_user gu	LEFT JOIN `user` u  ON u.`id` = gu.`user_id` Where gu.game_id=?").SetArgs(gid).QueryRows(&rmap)
	logs.Info(rmap)
	return rmap
}

func AddGameMatchLog(gml GameMatchLog) {
	o.Insert(&gml)
}

func GetUserPoint(uid int) int {
	u := User{
		Id: uid,
	}
	o.Read(&u)
	return u.Point
}

func ChangeUserPoint(uid int, point int) {
	u := User{
		Id: uid,
	}
	o.Read(&u)
	u.Point = u.Point + point
	o.Update(&u)
}
