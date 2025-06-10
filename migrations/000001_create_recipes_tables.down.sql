-- +goose Down
DROP TABLE IF EXISTS recipes_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipes; 