package models

type User struct {
	Id    int    `form:"-"`
	Name  string `form:"name"`
	Point string `form:"-"`
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
