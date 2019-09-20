package main

import (
	"./calc"
	"./model"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const helpText = "factorio_planner [name_of_product string] [target_quantity_per_hour int]"

func main() {
	goal, err := makeGoalFromArgs()
	if err != nil {
		panic(err)
	}
	recipeData, err := readRecipes()
	if err != nil {
		panic(err)
	}
	products, err := parseProducts(recipeData)
	if err != nil {
		panic(err)
	}
	err = sanityCheckProducts(products)
	if err != nil {
		panic(err)
	}
	err = calc.Goal(*goal, products)
	if err != nil {
		panic(err)
	}
}

func makeGoalFromArgs() (*model.Goal, error) {
	args := os.Args[1:]
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments:\n%v\n", helpText)
	}

	if len(args[0]) <= 0 {
		return nil, fmt.Errorf("product name cannot be empty:\n%v\n", helpText)
	}

	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("couldn't parse quantity argument to an int:\n%v\n", helpText)
	}

	goal := model.MakeGoal(args[0], quantity)
	return &goal, nil
}

func readRecipes() (string, error) {
	buf, err := ioutil.ReadFile("recipes.csv")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func parseProducts(data string) ([]model.Product, error) {
	var products []model.Product
	split := strings.Split(data, "\n")
	for i, s := range split {
		p, err := model.ParseProduct(s)
		if err != nil {
			return products, fmt.Errorf("parsing error on line %v: %v", i+1, err)
		}
		products = append(products, p)
	}
	return products, nil
}

func sanityCheckProducts(products []model.Product) error {
	for x, p1 := range products {
		if len(p1.Name) <= 0 {
			return fmt.Errorf("product name is empty on line %v", x+1)
		}
		if p1.BuildTimeSeconds < 0 {
			return fmt.Errorf("invalid build time for %v on line %v", p1.Name, x+1)
		}
		if p1.Output <= 0 {
			return fmt.Errorf("invalid output quantity for %v on line %v", p1.Name, x+1)
		}
		if len(p1.CreatedIn) <= 0 {
			return fmt.Errorf("created-in is empty for %v on line %v", p1.Name, x+1)
		}

		for y, p2 := range products {
			if x != y && p1.SanitizedName == p2.SanitizedName {
				return fmt.Errorf("found duplicate product of %v at line %v", p1.Name, y+1)
			}
		}

		for _, ing := range p1.Ingredients {
			var found = false
			for _, p2 := range products {
				if ing.SanitizedName == p2.SanitizedName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("unknown ingredient %v on line %v", ing.Name, x+1)
			}

			if ing.Quantity <= 0 {
				return fmt.Errorf("invalid ingredient quantity for %v on line %v", ing.Name, x+1)
			}
		}
	}
	return nil
}
