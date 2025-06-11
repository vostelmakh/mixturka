package repository

import (
	"context"
	"database/sql"

	"github.com/vostelmakh/mixturka/src/domain"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) SaveRecipe(ctx context.Context, recipe *domain.Recipe) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var recipeID int64
	err = tx.QueryRowContext(ctx,
		"INSERT INTO recipes (name) VALUES ($1) RETURNING id",
		recipe.Name,
	).Scan(&recipeID)
	if err != nil {
		return err
	}

	for _, ingredient := range recipe.Ingredients {
		var ingredientID int64
		err = tx.QueryRowContext(ctx,
			"INSERT INTO ingredients (recipe_id, name, quantity) VALUES ($1, $2, $3) RETURNING id",
			recipeID, ingredient.Name, ingredient.Quantity,
		).Scan(&ingredientID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx,
			"INSERT INTO recipes_ingredients (recipe_id, ingredient_id) VALUES ($1, $2)",
			recipeID, ingredientID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *RecipeRepository) GetRecipes(ctx context.Context) ([]domain.Recipe, error) {
	query := `
		SELECT r.id, r.name, i.id, i.name, i.quantity
		FROM recipes r
		LEFT JOIN recipes_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		ORDER BY r.id, i.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := make(map[int64]domain.Recipe)
	for rows.Next() {
		var recipeID int64
		var recipeName string
		var ingredientID sql.NullInt64
		var ingredientName sql.NullString
		var quantity sql.NullInt32

		err := rows.Scan(&recipeID, &recipeName, &ingredientID, &ingredientName, &quantity)
		if err != nil {
			return nil, err
		}

		recipe, exists := recipes[recipeID]
		if !exists {
			recipe = domain.Recipe{
				ID:          recipeID,
				Name:        recipeName,
				Ingredients: make([]domain.Ingredient, 0),
			}

			recipes[recipeID] = recipe
		}

		if ingredientID.Valid && ingredientName.Valid && quantity.Valid {
			recipe.Ingredients = append(recipe.Ingredients, domain.Ingredient{
				ID:       ingredientID.Int64,
				Name:     ingredientName.String,
				Quantity: int(quantity.Int32),
			})
		}
	}

	result := make([]domain.Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		result = append(result, recipe)
	}

	return result, nil
}
