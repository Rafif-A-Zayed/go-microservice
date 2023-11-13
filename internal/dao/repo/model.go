package repo

// UserModel represents table user.
type UserModel struct {
	ID string `gorm:"primarykey"`
	//CreatedAt  time.Time
	CustomerID string
	Status     string
	Address    string
}
