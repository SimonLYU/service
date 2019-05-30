package AccountModel

type Account struct {
	Name        string
	DatabaseName string
	AccountList []string
}

type SingleAccount struct {
	Name        string
	DatabaseName string
	Account      string
}

type GetAccount struct {
	DatabaseName string
}