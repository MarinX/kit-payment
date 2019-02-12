package account

// Service is the interface that provides account methods.
type Service interface {
	// CreateAccount creates new account with generated ID
	CreateAccount() (*Account, error)

	// GetAccount returns account by ID
	GetAccount(string) (*Account, error)

	// Accounts lists all accounts
	Accounts() []*Account

	// SetBalanceForAccount hard reset balance for account for given currency
	SetBalanceForAccount(*Account, Currency, float64) (*Account, error)
}

type service struct {
	accounts Repository
}

// NewService creates account service
func NewService(accounts Repository) Service {
	return &service{
		accounts: accounts,
	}
}

func (s *service) CreateAccount() (*Account, error) {
	acc := New()
	err := s.accounts.Store(acc)
	return acc, err
}

func (s *service) Accounts() []*Account {
	return s.accounts.FindAll()
}

func (s *service) GetAccount(id string) (*Account, error) {
	return s.accounts.Find(id)
}

func (s *service) SetBalanceForAccount(account *Account, currency Currency, amount float64) (*Account, error) {
	account.SetBalance(currency, amount)
	err := s.accounts.Store(account)
	return account, err
}
