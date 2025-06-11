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
