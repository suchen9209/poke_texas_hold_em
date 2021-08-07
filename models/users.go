package models

type User struct {
	Id    int    `form:"-"`
	Name  string `form:"name"`
	Point string `form:"-"`
}

func CheckUser(user *User) bool {
	can_login := false
	o.QueryTable("user").Filter("username", user.Name).One(user)
	if user.Id > 0 {
		can_login = true
	}
	return can_login

}
