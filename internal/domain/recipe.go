package domain

type Recipe struct {
	ID          int64        `db:"id"`
	Name        string       `db:"name"`
	Ingredients []Ingredient `db:"ingredients"`
}
