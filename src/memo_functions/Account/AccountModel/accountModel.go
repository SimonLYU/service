package AccountModel

type Account struct {
	Name        string
	DatabaseName string
	AccountList []string
}

type GetAccount struct {
	DatabaseName string
}