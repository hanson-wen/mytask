package model

type User struct {
	Uid  int
	Name string
}

// Auth parse the token check assert user info
func Auth(token string) (u User) {
	// todo check the token and auth
	u = User{
		Uid:  1,
		Name: "hanson",
	}
	return
}
