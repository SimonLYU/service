package UserModel

type User struct {
	Name 		string
	Account 	string
	Password 	string
	LinkInviteCode string
	LinkAccount string
}

type Password struct {
	Account 	string
	OldPassword string
	NewPassword	string
}