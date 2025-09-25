package recipe

import (
	"context"
	"encoding/json"

	"github.com/vostelmakh/mixturka/internal/domain"
	"github.com/vostelmakh/mixturka/internal/infrastructure/repository"
)

type Processor struct {
	repo *repository.RecipeRepository
}

func NewRecipeProcessor(repo *repository.RecipeRepository) *Processor {
	return &Processor{
		repo: repo,
	}
}

func (p *Processor) ProcessRecipe(ctx context.Context, message []byte) error {
	var recipe domain.Recipe
	if err := json.Unmarshal(message, &recipe); err != nil {
		return err
	}

	return p.repo.SaveRecipe(ctx, &recipe)
}

func (p *Processor) GetRecipes(ctx context.Context) ([]domain.Recipe, error) {
	recipes, err := p.repo.GetRecipes(ctx)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}
