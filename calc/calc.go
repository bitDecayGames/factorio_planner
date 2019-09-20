package calc

import (
	"../model"
	"fmt"
	"math"
	"sort"
)

func Goal(goal model.Goal, products []model.Product) error {
	ingredients, err := recursiveGoal(goal, products, 0)
	if err != nil {
		return err
	}
	combined := combineIngredients(ingredients)
	sortIngredients(combined)
	printIngredients(combined, products)

	return nil
}

func recursiveGoal(goal model.Goal, products []model.Product, depth int) ([]model.SubGoal, error) {
	if depth > 1000 {
		return nil, fmt.Errorf("infinite loop found with ingredient %v", goal.Name)
	}
	var rawMaterials []model.SubGoal
	product := findProduct(goal.Name, products)
	if product == nil {
		return nil, fmt.Errorf("failed to find product with name %v", goal.Name)
	}
	rawMaterials = append(rawMaterials, model.MakeSubGoal(product.Name, goal.QuantityPerHour, *product, depth))
	for _, ing := range product.Ingredients {
		subRawMaterials, err := recursiveGoal(model.MakeGoal(ing.Name, ing.Quantity * goal.QuantityPerHour), products, depth + 1)
		if err != nil {
			return nil, err
		}
		rawMaterials = append(rawMaterials, subRawMaterials...)
	}

	return rawMaterials, nil
}

func findProduct(name string, products []model.Product) *model.Product {
	sanitizedName := model.SanitizeName(name)
	for i, p := range products {
		if p.SanitizedName == sanitizedName {
			return &products[i]
		}
	}
	return nil
}

func combineIngredients(ingredients []model.SubGoal) []model.SubGoal {
	var combined []model.SubGoal
	for x, ing1 := range ingredients {
		var alreadyCombined = false
		for _, ing2 := range combined {
			if ing1.SanitizedName == ing2.SanitizedName {
				alreadyCombined = true
				break
			}
		}
		if !alreadyCombined {
			combinedIngredient := model.MakeSubGoal(ing1.Name, ing1.QuantityPerHour, ing1.Product, ing1.Depth)
			for z, ing3 := range ingredients {
				if z > x && combinedIngredient.SanitizedName == ing3.SanitizedName {
					combinedIngredient.QuantityPerHour += ing3.QuantityPerHour
					if combinedIngredient.Depth < ing3.Depth {
						combinedIngredient.Depth = ing3.Depth
					}
				}
			}
			combined = append(combined, combinedIngredient)
		}
	}
	return combined
}

func sortIngredients(ingredients []model.SubGoal) {
	sort.SliceStable(ingredients, func(a, b int) bool {
		var aI = ingredients[a]
		var bI = ingredients[b]

		if aI.Depth == bI.Depth {
			return aI.QuantityPerHour <= bI.QuantityPerHour
		}
		return aI.Depth <= bI.Depth
	})
}

func printIngredients(ingredients []model.SubGoal, products []model.Product) {
	for _, ing := range ingredients {
		var product = findProduct(ing.Name, products)
		if product != nil {
			var count int
			if ing.Product.BuildTimeSeconds <= 0 {
				count = 1
			} else {
				count = int(math.Ceil(float64(ing.QuantityPerHour) / 60.0 / 60.0 * ing.Product.BuildTimeSeconds))
			}
			var createdIn = ing.Product.CreatedIn
			if count > 1 {
				createdIn = fmt.Sprintf("%vs", createdIn)
			}
			fmt.Printf("%v: %v in %v %v\n", ing.Name, ing.QuantityPerHour, count, createdIn)
		}
	}
}