package processor

import (
	"context"
	"encoding/json"

	"github.com/gbrayhan/microservices-go/src/domain"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"
)

type RecipeProcessor struct {
	repo *repository.RecipeRepository
}

func NewRecipeProcessor(repo *repository.RecipeRepository) *RecipeProcessor {
	return &RecipeProcessor{
		repo: repo,
	}
}

func (p *RecipeProcessor) ProcessRecipe(ctx context.Context, message []byte) error {
	var recipe domain.Recipe
	if err := json.Unmarshal(message, &recipe); err != nil {
		return err
	}

	return p.repo.SaveRecipe(ctx, &recipe)
}
