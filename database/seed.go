package database

import "cosmetics_catalog/models"

func SeedTestData() error {
	// Создание брендов
	brands := []models.Brand{
		{Name: "L'Oreal", Slug: "loreal"},
		{Name: "Maybelline", Slug: "maybelline"},
		{Name: "Nivea", Slug: "nivea"},
	}

	for _, brand := range brands {
		if err := DB.Create(&brand).Error; err != nil {
			return err
		}
	}

	// Создание категорий и подкатегорий
	categories := []models.Category{
		{
			Name: "Макияж",
			Slug: "makiyazh",
			Subcategories: []models.Subcategory{
				{Name: "Лицо", Slug: "litso"},
				{Name: "Глаза", Slug: "glaza"},
				{Name: "Губы", Slug: "guby"},
			},
		},
		{
			Name: "Уход",
			Slug: "ukhod",
			Subcategories: []models.Subcategory{
				{Name: "Очищение", Slug: "ochishenie"},
				{Name: "Увлажнение", Slug: "uvlazhnenie"},
				{Name: "Маски", Slug: "maski"},
			},
		},
	}

	for _, category := range categories {
		if err := DB.Create(&category).Error; err != nil {
			return err
		}
	}

	// Создание тестовых продуктов
	products := []models.Product{
		{
			Name:          "Тональный крем",
			Slug:          "tonalnyi-krem",
			BrandID:       1,
			SubcategoryID: 1, // Лицо
			Price:         1299.99,
			ImagePath:     "images/creme.jpg",
			Description:   "Легкий тональный крем с натуральным покрытием",
		},
		{
			Name:          "Тушь для ресниц",
			Slug:          "tush-dlya-resnits",
			BrandID:       2,
			SubcategoryID: 2, // Глаза
			Price:         899.99,
			ImagePath:     "images/mascara.jpg",
			Description:   "Объемная тушь для ресниц",
			IsOnSale:      true,
			SalePrice:     699.99,
		},
	}

	for _, product := range products {
		if err := DB.Create(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
