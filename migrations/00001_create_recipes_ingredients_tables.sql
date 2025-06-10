-- +goose Up
CREATE TABLE recipes (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE ingredients (
    id BIGSERIAL PRIMARY KEY,
    recipe_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    CONSTRAINT fk_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes (id)
);

CREATE TABLE recipes_ingredients (
    id BIGSERIAL PRIMARY KEY,
    recipe_id BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    CONSTRAINT fk_recipe_id FOREIGN KEY (recipe_id) REFERENCES recipes(id),
    CONSTRAINT fk_ingredient_id FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);

CREATE INDEX idx_recipe_id ON recipes_ingredients(recipe_id);
CREATE INDEX idx_ingredient_id ON recipes_ingredients(ingredient_id);

-- +goose Down
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipes_ingredients;