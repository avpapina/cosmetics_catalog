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
			Name:          "Увлажняющий крем",
			Slug:          "uvlazhnyayushchii-krem",
			BrandID:       1,
			SubcategoryID: 1, // Лицо
			Price:         1499.00,
			ImagePath:     "images/moisturizer.jpg",
			Description:   "Интенсивное увлажнение на 24 часа",
			IsOnSale:      true,
			SalePrice:     1199.00,
		},
		{
			Name:          "Очищающая пенка",
			Slug:          "ochishchayushchaya-penka",
			BrandID:       3,
			SubcategoryID: 1, // Лицо
			Price:         799.50,
			ImagePath:     "images/cleanser.jpg",
			Description:   "Мягкое очищение без стягивания",
		},
		{
			Name:          "Сыворотка с гиалуроновой кислотой",
			Slug:          "syvorotka-s-gialuronovoi-kislotoy",
			BrandID:       2,
			SubcategoryID: 1, // Лицо
			Price:         2499.00,
			ImagePath:     "images/serum.jpg",
			Description:   "Глубокое увлажнение и разглаживание морщин",
		},
		{
			Name:          "BB-крем",
			Slug:          "bb-krem",
			BrandID:       4,
			SubcategoryID: 1, // Лицо
			Price:         1599.00,
			ImagePath:     "images/bb-cream.jpg",
			Description:   "Многофункциональный уход и макияж в одном",
			IsOnSale:      true,
			SalePrice:     1299.00,
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
