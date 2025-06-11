package brew

import (
	"context"
	"fmt"

	"github.com/vostelmakh/mixturka/src/domain"
	"github.com/vostelmakh/mixturka/src/infrastructure/repository"
)

const (
	successfulBrew = true
	failedBrew     = false
)

type Ingredient struct {
	Name     string
	Quantity int
}

type Processor struct {
	repo *repository.RecipeRepository
}

func NewGRPCProcessor(repo *repository.RecipeRepository) *Processor {
	return &Processor{
		repo: repo,
	}
}

func (p *Processor) BrewPot(ctx context.Context, ingredients []Ingredient) (bool, error) {
	recipesList, err := p.repo.GetRecipes(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get recipes: %w", err)
	}

	brewIngredients := make(map[string]int)
	for _, ingredient := range ingredients {
		brewIngredients[ingredient.Name] = ingredient.Quantity
	}

	for _, recipe := range recipesList {
		if p.canBrew(brewIngredients, recipe.Ingredients) {
			fmt.Printf("Successfully brewed %s!\n", recipe.Name)

			return successfulBrew, nil
		}
	}

	return failedBrew, nil
}
func (p *Processor) canBrew(brewIngredients map[string]int, recipeIngredients []domain.Ingredient) bool {
	recipeIngredientsMap := make(map[string]int)
	for _, ingredient := range recipeIngredients {
		recipeIngredientsMap[ingredient.Name] = ingredient.Quantity
	}

	for brewIngredientName, brewIngredientQuantity := range brewIngredients {
		ingredientQuantity, ok := recipeIngredientsMap[brewIngredientName]
		if !ok || brewIngredientQuantity > ingredientQuantity {
			return false
		}
	}

	return true
}
