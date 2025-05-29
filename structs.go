package main

//структура продукта
type Product struct {
	Name  string
	Price float64
}

//структура категории в каталоге
type Category struct {
	Name     string        `json:"name"`     // Название категории
	Slug     string        `json:"slug"`     // Часть URL
	Children []SubCategory `json:"children"` // Подкатегории
}

//Структура для подкатегорий
type SubCategory struct {
	Name string `json:"name"` // Название категории
	Slug string `json:"slug"` // Часть URL
}

//удалить
var products = []Product{
	{Name: "milk", Price: 89.99},
	{Name: "bread", Price: 45.50},
	{Name: "apples", Price: 129.99},
	{Name: "coffee", Price: 349.90},
	{Name: "chocolate", Price: 79.75},
}

var brands = []SubCategory{
	{Name: "L'Oreal", Slug: "loreal"},
	{Name: "Maybelline", Slug: "maybelline"},
}

var makeup = []SubCategory{
	{Name: "Лицо", Slug: "face"},
	{Name: "Глаза", Slug: "eyes"},
	{Name: "Губы", Slug: "lips"},
	{Name: "Брови", Slug: "brows"},
}

var care = []SubCategory{
	{Name: "Очищение", Slug: "cleansing"},
	{Name: "Тонизирование", Slug: "toning"},
}

// Основной каталог
var catalog = []Category{
	{
		Name:     "Бренды",
		Slug:     "brands",
		Children: brands,
	},
	{
		Name:     "Акции",
		Slug:     "sales",
		Children: []SubCategory{},
	},
	{
		Name:     "Макияж",
		Slug:     "makeup",
		Children: makeup,
	},
	{
		Name:     "Уход",
		Slug:     "care",
		Children: care,
	},
}
