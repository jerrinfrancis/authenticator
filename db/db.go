package db

type UserSecret struct {
	Passwd string
	Salt   string
	Secret string
}

type UserInfo struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber uint64
}

type User struct {
	UserInfo
	UserSecret
}

type UserDB interface {
	Insert(u User) (*User, error)
	Find(email string) (*User, error)
}
type DB interface {
	User() UserDB
}
