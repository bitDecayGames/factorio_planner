package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Product struct {
	Name string
	BuildTimeSeconds float64
	Output int
	CreatedIn string
	Ingredients []Ingredient

	SanitizedName string
}

func parseProduct(data string) (Product, error) {
	split := strings.Split(data, ",")
	if len(split) != 5 {
		return Product{}, fmt.Errorf("must have 5 items separated by ','")
	}
	output, err := strconv.Atoi(split[2])
	if err != nil {
		return Product{}, err
	}
	buildTime, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return Product{}, err
	}
	ing, err := parseIngredients(split[4])
	if err != nil {
		return Product{}, err
	}
	return Product{
		Name: strings.TrimSpace(split[0]),
		SanitizedName: sanitizeName(split[0]),
		BuildTimeSeconds: buildTime,
		Output: output,
		CreatedIn: strings.TrimSpace(split[3]),
		Ingredients: ing,
	}, nil
}

type Ingredient struct {
	Name string
	Quantity int

	SanitizedName string
}

func parseIngredients(data string) ([]Ingredient, error) {
	var ingredients []Ingredient
	if len(strings.TrimSpace(data)) <= 0 {
		return ingredients, nil
	}
	list := strings.Split(data, "|")
	for _, s := range list {
		split := strings.Split(s, ":")
		if len(split) != 2 {
			return ingredients, fmt.Errorf("each ingredient must be in the format <name>:<quantity>")
		}
		quantity, err := strconv.Atoi(split[1])
		if err != nil {
			return ingredients, err
		}
		ingredients = append(ingredients, Ingredient{
			Name: split[0],
			SanitizedName: sanitizeName(split[0]),
			Quantity: quantity,
		})
	}
	return ingredients, nil
}

type Goal struct {
	Name string
	QuantityPerMinute int

	SanitizedName string
}

func MakeGoal(name string, quantityPerMinute int) Goal {
	return Goal {
		Name: name,
		SanitizedName: sanitizeName(name),
		QuantityPerMinute: quantityPerMinute,
	}
}

func sanitizeName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "")
}