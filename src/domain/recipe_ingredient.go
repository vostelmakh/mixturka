package domain

type RecipeIngredient struct {
	ID           int64 `db:"id"`
	RecipeID     int64 `db:"recipe_id"`
	IngredientID int64 `db:"ingredient_id"`
}
