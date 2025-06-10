package domain

type Ingredient struct {
	ID       int64  `db:"id"`
	RecipeID int64  `db:"recipe_id"`
	Name     string `db:"name"`
	Quantity int    `db:"quantity"`
}
