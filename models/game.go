package models

import "github.com/beego/beego/v2/core/logs"

// Game 一整场游戏
type Game struct {
	Id       int    `form:"-"`
	GameName string `form:"name"`
}

//GameMatch 一局游戏，从发牌到结算
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

// GameMatchLog 单局游戏玩家操作记录
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

// GameUser 一整场游戏玩家参与位置，是否在线
type GameUser struct {
	Id       int
	UserId   int
	GameId   int	//即RoomID
	Position int
	Online   int
}

type UserPointSeat struct {
	UserId   int
	Name     string
	Position int
	Point    int
}

type EndGameCardInfo struct {
	Type       EventType
	WinPos     []int
	WinCard    [][]Card
	PublicCard []Card
	BigCard    []Card
	UserCards  map[int][]Card
}

const (
	GAME_OP_RAISE  = "raise"
	GAME_OP_CALL   = "call"
	GAME_OP_CHECK  = "check"
	GAME_OP_FOLD   = "fold"
	GAME_OP_ALLIN  = "allin"
	GAMEUSERNUMBER = 8

	GAME_STATUS_END       = "END"
	GAME_STATUS_ROUND1    = "ROUND1"
	GAME_STATUS_ROUND2    = "ROUND2"
	GAME_STATUS_ROUND3    = "ROUND3"
	GAME_STATUS_ROUND4    = "ROUND4"
	GAME_STATUS_LICENSING = "LICENSING"
)

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

func SetUserIntoRoom(u User,RoomID int) GameUser {
	//如果已经存在，则直接返回位置
	var gut = &GameUser{
		UserId: u.Id,
		GameId: RoomID,
	}
	o.Read(gut, "UserId", "GameId")
	if gut.Id > 0 {
		gut.Online = 1
		o.Update(gut)
		return *gut
	}

	//如果不存在查找剩余位置，分配座位
	var gu []*GameUser
	number, _ := o.QueryTable("game_user").Filter("game_id", RoomID).Filter("online", 1).All(&gu)
	if number >= GAMEUSERNUMBER {
		panic("too many player")
	}
	var positionList = make(map[int]int)
	for _, g := range gu {
		positionList[g.Position] = g.UserId
	}
	var tmpG GameUser
	var pos int
	for i := 1; i <= GAMEUSERNUMBER; i++ {
		if positionList[i] == 0 {
			pos = i
			tmpG.UserId = u.Id
			tmpG.GameId = RoomID
			tmpG.Position = pos
			tmpG.Online = 1
			o.Insert(&tmpG)
			break
		}
	}

	return tmpG
}

func GetRoomUserPositionList(roomId int) map[int]int{
	var gu []*GameUser
	_, _ = o.QueryTable("game_user").Filter("game_id", roomId).Filter("online", 1).All(&gu)
	var positionList = make(map[int]int)
	for _, g := range gu {
		positionList[g.Position] = g.UserId
	}
	return positionList
}


func RemoveGameUser(uid int, gid int) {
	o.QueryTable("game_user").Filter("game_id", gid).Filter("user_id", uid).Delete()
}

func TruncateGameUser() {
	o.Raw("TRUNCATE game_user").Exec()
}

func GetUserPointWithSeat(gid int) []UserPointSeat {
	var gu []UserPointSeat
	o.Raw(`select gu.user_id,u.name,gu.position, u.point from game_user as gu left join user as u on u.id = gu.user_id where gu.game_id = ?`).SetArgs(1).QueryRows(&gu)
	return gu
}

func InitGameMatch(gid int) GameMatch {
	var tmp, newGame GameMatch
	var gameUserList []GameUser
	o.QueryTable("game_user").Filter("game_id", gid).Filter("online", 1).All(&gameUserList)
	if gameUserList == nil {
		panic("no user")
	}
	tmpMap := make(map[int]int)
	for _, v := range gameUserList {
		tmpMap[v.Position] = v.UserId
	}
	newGame.GameId = gid
	o.QueryTable("game_match").Filter("game_id", gid).OrderBy("-id").One(&tmp)
	if tmp.Id <= 0 {
		//之前无对局
		newGame.SmallBindPosition = getPosition(0, tmpMap)
	} else {
		newGame.SmallBindPosition = getPosition(tmp.SmallBindPosition, tmpMap)
	}
	newGame.BigBindPosition = getPosition(newGame.SmallBindPosition, tmpMap)
	newGame.GameStatus = "INIT"
	o.Insert(&newGame)
	return newGame

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
	//logs.Info(rmap)
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
