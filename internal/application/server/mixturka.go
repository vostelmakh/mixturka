package server

import (
	"context"

	"github.com/vostelmakh/mixturka/internal/application/processor/brew"
	"github.com/vostelmakh/mixturka/internal/application/processor/recipe"
	mixturkaGrpc "github.com/vostelmakh/mixturka/internal/infrastructure/grpc"
)

type MixturkaServer struct {
	mixturkaGrpc.UnimplementedMixturkaServer
	recipeProcessor *recipe.Processor
	brewProcessor   *brew.Processor
}

func NewMixturkaServer(recipeProcessor *recipe.Processor, brewProcessor *brew.Processor) *MixturkaServer {
	return &MixturkaServer{
		recipeProcessor: recipeProcessor,
		brewProcessor:   brewProcessor,
	}
}

func (s *MixturkaServer) GetRecipes(ctx context.Context, req *mixturkaGrpc.GetRecipesRequest) (*mixturkaGrpc.GetRecipesResponse, error) {
	recipes, err := s.recipeProcessor.GetRecipes(ctx)
	if err != nil {
		return nil, err
	}

	response := &mixturkaGrpc.GetRecipesResponse{
		Recipes: make([]*mixturkaGrpc.Recipe, 0, len(recipes)),
	}

	for _, recipe := range recipes {
		grpcRecipe := &mixturkaGrpc.Recipe{
			Id:          recipe.ID,
			Name:        recipe.Name,
			Ingredients: make([]*mixturkaGrpc.Ingredient, 0, len(recipe.Ingredients)),
		}

		for _, ingredient := range recipe.Ingredients {
			grpcRecipe.Ingredients = append(grpcRecipe.Ingredients, &mixturkaGrpc.Ingredient{
				Id:       ingredient.ID,
				Name:     ingredient.Name,
				Quantity: int32(ingredient.Quantity),
			})
		}

		response.Recipes = append(response.Recipes, grpcRecipe)
	}

	return response, nil
}

func (s *MixturkaServer) BrewPot(ctx context.Context, req *mixturkaGrpc.PotBrewRequest) (*mixturkaGrpc.PotBrewResponse, error) {
	// Преобразуем ингредиенты из gRPC в доменные модели
	ingredients := make([]brew.Ingredient, 0, len(req.Ingredients))
	for _, ing := range req.Ingredients {
		ingredients = append(ingredients, brew.Ingredient{
			Name:     ing.Name,
			Quantity: int(ing.Quantity),
		})
	}

	// Запускаем процесс варки
	started, err := s.brewProcessor.BrewPot(ctx, ingredients)
	if err != nil {
		return &mixturkaGrpc.PotBrewResponse{
			Started: false,
			Error: &mixturkaGrpc.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	return &mixturkaGrpc.PotBrewResponse{
		Started: started,
	}, nil
}
