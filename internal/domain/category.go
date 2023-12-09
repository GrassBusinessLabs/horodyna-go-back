package domain

type Category string

const (
	VegetablesCategory        Category = "Vegetables"
	FruitsCategory            Category = "Fruits"
	MilkProductsCategory      Category = "Milk Products"
	FishCategory              Category = "Fish"
	MeatCategory              Category = "Meat"
	BakeryCategory            Category = "Bakery"
	FrozenFoodsCategory       Category = "Frozen Foods"
	SweetsCategory            Category = "Sweets"
	HealthAndWellnessCategory Category = "Health and Wellness"
)

func GetCategoriesList() []Category {
	return []Category{
		VegetablesCategory,
		FishCategory,
		FrozenFoodsCategory,
		FruitsCategory,
		BakeryCategory,
		SweetsCategory,
		HealthAndWellnessCategory,
		MeatCategory,
		MilkProductsCategory,
	}
}
