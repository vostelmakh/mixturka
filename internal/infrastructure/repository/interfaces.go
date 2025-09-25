package repository

import (
	"context"

	"github.com/vostelmakh/mixturka/internal/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_repository.go

type RecipeRepositoryInterface interface {
	GetRecipes(ctx context.Context) ([]domain.Recipe, error)
	SaveRecipe(ctx context.Context, recipe *domain.Recipe) error
}
