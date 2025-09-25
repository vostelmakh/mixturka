package brew

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vostelmakh/mixturka/internal/domain"
	mock_repository "github.com/vostelmakh/mixturka/internal/infrastructure/repository/mocks"
)

func TestProcessor_BrewPot(t *testing.T) {
	tests := []struct {
		name               string
		ingredients        []Ingredient
		mockSetup          func(*mock_repository.MockRecipeRepositoryInterface)
		expectedResult     bool
		expectedError      string
		expectedPrintCount int
	}{
		{
			name: "успешное варение - точное соответствие ингредиентов",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{
						{
							ID:   1,
							Name: "Торт",
							Ingredients: []domain.Ingredient{
								{Name: "мука", Quantity: 100},
								{Name: "сахар", Quantity: 50},
							},
						},
					}, nil)
			},
			expectedResult: true,
			expectedError:  "",
		},
		{
			name: "успешное варение - ингредиентов больше чем нужно",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 80},
				{Name: "сахар", Quantity: 30},
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{
						{
							ID:   1,
							Name: "Печенье",
							Ingredients: []domain.Ingredient{
								{Name: "мука", Quantity: 100},
								{Name: "сахар", Quantity: 50},
							},
						},
					}, nil)
			},
			expectedResult: true,
			expectedError:  "",
		},
		{
			name: "неуспешное варение - недостаточно ингредиентов",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 150},
				{Name: "сахар", Quantity: 30},
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{
						{
							ID:   1,
							Name: "Торт",
							Ingredients: []domain.Ingredient{
								{Name: "мука", Quantity: 100},
								{Name: "сахар", Quantity: 50},
							},
						},
					}, nil)
			},
			expectedResult: false,
			expectedError:  "",
		},
		{
			name: "неуспешное варение - отсутствует необходимый ингредиент",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "перец", Quantity: 10}, // этого ингредиента нет в рецепте
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{
						{
							ID:   1,
							Name: "Торт",
							Ingredients: []domain.Ingredient{
								{Name: "мука", Quantity: 100},
								{Name: "сахар", Quantity: 50},
							},
						},
					}, nil)
			},
			expectedResult: false,
			expectedError:  "",
		},
		{
			name: "неуспешное варение - пустой список рецептов",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{}, nil)
			},
			expectedResult: false,
			expectedError:  "",
		},
		{
			name: "ошибка при получении рецептов",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 100},
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return(nil, errors.New("ошибка базы данных"))
			},
			expectedResult: false,
			expectedError:  "failed to get recipes",
		},
		{
			name: "успешное варение - несколько рецептов, подходит первый",
			ingredients: []Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{
						{
							ID:   1,
							Name: "Торт",
							Ingredients: []domain.Ingredient{
								{Name: "мука", Quantity: 100},
								{Name: "сахар", Quantity: 50},
							},
						},
						{
							ID:   2,
							Name: "Хлеб",
							Ingredients: []domain.Ingredient{
								{Name: "мука", Quantity: 200},
								{Name: "дрожжи", Quantity: 10},
							},
						},
					}, nil)
			},
			expectedResult: true,
			expectedError:  "",
		},
		{
			name:        "успешное варение - пустой список ингредиентов для варки",
			ingredients: []Ingredient{},
			mockSetup: func(mockRepo *mock_repository.MockRecipeRepositoryInterface) {
				mockRepo.EXPECT().
					GetRecipes(gomock.Any()).
					Return([]domain.Recipe{
						{
							ID:          1,
							Name:        "Пустой рецепт",
							Ingredients: []domain.Ingredient{},
						},
					}, nil)
			},
			expectedResult: true,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockRecipeRepositoryInterface(ctrl)
			tt.mockSetup(mockRepo)

			processor := NewGRPCProcessor(mockRepo)
			ctx := context.Background()

			// Act
			result, err := processor.BrewPot(ctx, tt.ingredients)

			// Assert
			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestProcessor_canBrew тестирует внутреннюю логику сравнения ингредиентов
func TestProcessor_canBrew(t *testing.T) {
	tests := []struct {
		name              string
		brewIngredients   map[string]int
		recipeIngredients []domain.Ingredient
		expected          bool
	}{
		{
			name: "точное соответствие",
			brewIngredients: map[string]int{
				"мука":  100,
				"сахар": 50,
			},
			recipeIngredients: []domain.Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			expected: true,
		},
		{
			name: "больше ингредиентов чем нужно",
			brewIngredients: map[string]int{
				"мука":  80,
				"сахар": 30,
			},
			recipeIngredients: []domain.Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			expected: true,
		},
		{
			name: "недостаточно одного ингредиента",
			brewIngredients: map[string]int{
				"мука":  150,
				"сахар": 30,
			},
			recipeIngredients: []domain.Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			expected: false,
		},
		{
			name: "ингредиент отсутствует в рецепте",
			brewIngredients: map[string]int{
				"мука":  100,
				"перец": 10, // этого ингредиента нет в рецепте
			},
			recipeIngredients: []domain.Ingredient{
				{Name: "мука", Quantity: 100},
				{Name: "сахар", Quantity: 50},
			},
			expected: false,
		},
		{
			name:            "пустые ингредиенты для варки - можно сварить любой рецепт",
			brewIngredients: map[string]int{},
			recipeIngredients: []domain.Ingredient{
				{Name: "мука", Quantity: 100},
			},
			expected: true, // пустой список ингредиентов означает что ничего не требуется
		},
		{
			name: "пустой рецепт - нельзя использовать ингредиенты",
			brewIngredients: map[string]int{
				"мука": 100,
			},
			recipeIngredients: []domain.Ingredient{},
			expected:          false, // в рецепте нет ингредиентов, поэтому мука не подходит
		},
		{
			name:              "пустой рецепт и пустые ингредиенты",
			brewIngredients:   map[string]int{},
			recipeIngredients: []domain.Ingredient{},
			expected:          true, // ничего не нужно для пустого рецепта
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			processor := &Processor{}

			// Act
			result := processor.canBrew(tt.brewIngredients, tt.recipeIngredients)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
