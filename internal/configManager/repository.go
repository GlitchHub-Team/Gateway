package configmanager

type SQLiteConfigRepository struct{}

func NewSQLiteConfigRepository() *SQLiteConfigRepository {
	return &SQLiteConfigRepository{}
}
