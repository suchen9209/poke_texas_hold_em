package models

type User struct {
	Id       int    `form:"-"`
	Name     string `form:"name"`
	Point    int    `form:"-"`
	Password string ``
}

func CheckUserLogin(user *User) error {
	var e error
	if len(user.Password) > 0 && len(user.Name) > 0 {
		e = o.Read(user, "Name", "Password")
	}
	return e
}

func AddUser(user *User) error {
	_, err := o.Insert(user)
	return err
}

func CheckUser(name string) (*User, bool) {
	can_login := false
	var user User
	o.QueryTable("user").Filter("name", name).One(&user)
	if user.Id > 0 {
		can_login = true
	}
	return &user, can_login

}

func CheckUserInGame(name string) (*GameUser, bool) {
	can_login := false
	var user User
	var gameUser GameUser
	o.QueryTable("user").Filter("name", name).One(&user)
	if user.Id > 0 {
		o.QueryTable("game_user").Filter("user_id", user.Id).One(&gameUser)
		if gameUser.Id > 0 {
			return &gameUser, true
		}
	}
	return &gameUser, can_login
}
