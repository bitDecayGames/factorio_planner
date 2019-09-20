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

func ParseProduct(data string) (Product, error) {
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
	ing, err := ParseIngredients(split[4])
	if err != nil {
		return Product{}, err
	}
	return Product{
		Name: strings.TrimSpace(split[0]),
		SanitizedName: SanitizeName(split[0]),
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

func MakeIngredient(name string, quantity int) Ingredient {
	return Ingredient{
		Name:          name,
		Quantity:      quantity,
		SanitizedName: SanitizeName(name),
	}
}

func ParseIngredients(data string) ([]Ingredient, error) {
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
		ingredients = append(ingredients, MakeIngredient(split[0], quantity))
	}
	return ingredients, nil
}

type Goal struct {
	Name string
	QuantityPerHour int

	SanitizedName string
}

func MakeGoal(name string, quantityPerHour int) Goal {
	return Goal {
		Name: name,
		SanitizedName: SanitizeName(name),
		QuantityPerHour: quantityPerHour,
	}
}

type SubGoal struct {
	Name string
	QuantityPerHour int
	SanitizedName string
	Product Product
	Depth int
}

func MakeSubGoal(name string, quantityPerHour int, product Product, depth int) SubGoal {
	return SubGoal {
		Name: name,
		SanitizedName: SanitizeName(name),
		QuantityPerHour: quantityPerHour,
		Product: product,
		Depth: depth,
	}
}

func SanitizeName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "")
}