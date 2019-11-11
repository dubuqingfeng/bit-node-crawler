package models

type Result struct {
	Status    string
	Address   string
	Peers     string
	UserAgent string
	ErrMsg    string
	CoinType  string
	Timestamp string
	Height    int64
}
