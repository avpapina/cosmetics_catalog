package database

import (
	"cosmetics_catalog/models"
)

func SeedTestData() error {
	// 1. Бренды
	brands := []models.Brand{
		{Name: "Bioderma", Slug: "bioderma"},
		{Name: "L'Oréal Paris", Slug: "loreal-paris"},
		{Name: "Pusy", Slug: "pusy"},
		{Name: "Dior", Slug: "dior"},
		{Name: "Dr. Jart+", Slug: "dr-jart"},
		{Name: "Clarins", Slug: "clarins"},
		{Name: "Estée Lauder", Slug: "estee-lauder"},
		{Name: "Catrice", Slug: "catrice"},
		{Name: "Clinique", Slug: "clinique"},
		{Name: "Shik", Slug: "shik"},
		{Name: "Kiko Milano", Slug: "kiko-milano"},
		{Name: "Erborian", Slug: "erborian"},
	}
	for _, brand := range brands {
		if err := DB.Create(&brand).Error; err != nil {
			return err
		}
	}

	// 2. Категории с подкатегориями
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
				{Name: "Тонизирование", Slug: "tonizirovanie"},
			},
		},
	}
	for _, category := range categories {
		if err := DB.Create(&category).Error; err != nil {
			return err
		}
	}

	// 3. Получение нужных подкатегорий из БД
	var (
		makeupFace, makeupEyes, makeupLips       models.Subcategory
		careCleansing, careHydration, careToning models.Subcategory
	)
	subMap := map[string]*models.Subcategory{
		"Лицо":          &makeupFace,
		"Глаза":         &makeupEyes,
		"Губы":          &makeupLips,
		"Очищение":      &careCleansing,
		"Увлажнение":    &careHydration,
		"Тонизирование": &careToning,
	}
	for name, ptr := range subMap {
		if err := DB.Where("name = ?", name).First(ptr).Error; err != nil {
			return err
		}
	}

	// 4. Продукты
	products := []models.Product{
		{
			Name:          "Тональный крем",
			Slug:          "tonalnyi-krem",
			BrandID:       1,
			SubcategoryID: makeupFace.ID,
			Price:         1299.99,
			ImagePath:     "images/creme.jpg",
			Description:   "Легкий тональный крем с натуральным покрытием",
		},
		{
			Name:          "Увлажняющий крем",
			Slug:          "uvlazhnyayushchii-krem",
			BrandID:       1,
			SubcategoryID: careHydration.ID,
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
			SubcategoryID: careCleansing.ID,
			Price:         799.50,
			ImagePath:     "images/cleanser.jpg",
			Description:   "Мягкое очищение без стягивания",
		},
		{
			Name:          "Сыворотка с гиалуроновой кислотой",
			Slug:          "syvorotka-s-gialuronovoi-kislotoy",
			BrandID:       2,
			SubcategoryID: careHydration.ID,
			Price:         2499.00,
			ImagePath:     "images/serum.jpg",
			Description:   "Глубокое увлажнение и разглаживание морщин",
		},
		{
			Name:          "BB-крем",
			Slug:          "bb-krem",
			BrandID:       4,
			SubcategoryID: makeupFace.ID,
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
			SubcategoryID: makeupEyes.ID,
			Price:         899.99,
			ImagePath:     "images/mascara.jpg",
			Description:   "Объемная тушь для ресниц",
			IsOnSale:      true,
			SalePrice:     699.99,
		},
		{
			Name:          "Бальзам для губ",
			Slug:          "balzam-dlya-gub",
			BrandID:       5,
			SubcategoryID: makeupLips.ID,
			Price:         499.00,
			ImagePath:     "images/lip-balm.jpg",
			Description:   "Увлажняющий бальзам для губ",
		},
		{
			Name:          "Тоник для лица",
			Slug:          "tonik-dlya-litsa",
			BrandID:       6,
			SubcategoryID: careToning.ID,
			Price:         899.00,
			ImagePath:     "images/toner.jpg",
			Description:   "Тонизирующее средство для свежести кожи",
		},
		// bioderma
		{
			Name:          "Гель для умывания Sensibio Foaming Gel",
			Slug:          "bioderma-sensibio-foaming-gel",
			BrandID:       1,
			SubcategoryID: careCleansing.ID,
			Price:         1450.00,
			ImagePath:     "images/bioderma/sensibio-foaming-gel.jpg",
			Description:   "Мягкий пенящийся гель для чувствительной кожи",
			IsOnSale:      false,
		},
		{
			Name:          "Мицеллярная вода Sensibio H2O",
			Slug:          "bioderma-sensibio-h2o",
			BrandID:       1,
			SubcategoryID: careCleansing.ID,
			Price:         1290.00,
			ImagePath:     "images/bioderma/sensibio-h2o.jpg",
			Description:   "Легендарная мицеллярная вода для чувствительной кожи",
			IsOnSale:      true,
			SalePrice:     1099.00,
		},
		// L'Oréal Paris
		{
			Name:          "Тональный крем True Match",
			Slug:          "loreal-true-match",
			BrandID:       2,
			SubcategoryID: makeupFace.ID,
			Price:         899.00,
			ImagePath:     "images/loreal/true-match.jpg",
			Description:   "Тональный крем с естественным покрытием",
			IsOnSale:      false,
		},
		{
			Name:          "Тушь для ресниц Lash Paradise",
			Slug:          "loreal-lash-paradise",
			BrandID:       2,
			SubcategoryID: makeupEyes.ID,
			Price:         799.00,
			ImagePath:     "images/loreal/lash-paradise.jpg",
			Description:   "Объемная тушь для эффекта накладных ресниц",
			IsOnSale:      true,
			SalePrice:     699.00,
		},

		// Dior
		{
			Name:          "Помада Dior Addict Lip Glow",
			Slug:          "dior-lip-glow",
			BrandID:       4,
			SubcategoryID: makeupLips.ID,
			Price:         3200.00,
			ImagePath:     "images/dior/lip-glow.jpg",
			Description:   "Бальзам для губ с эффектом сияния",
			IsOnSale:      false,
		},
		{
			Name:          "Тушь для ресниц Diorshow",
			Slug:          "dior-diorshow",
			BrandID:       4,
			SubcategoryID: makeupEyes.ID,
			Price:         2900.00,
			ImagePath:     "images/dior/diorshow.jpg",
			Description:   "Культовая тушь для ресниц",
			IsOnSale:      true,
			SalePrice:     2500.00,
		},

		// Dr. Jart+
		{
			Name:          "BB-крем Premium Beauty Balm",
			Slug:          "drjart-bb-cream",
			BrandID:       5,
			SubcategoryID: makeupFace.ID,
			Price:         2800.00,
			ImagePath:     "images/drjart/bb-cream.jpg",
			Description:   "Многофункциональный BB-крем с уходом",
			IsOnSale:      false,
		},
	}
	for _, product := range products {
		if err := DB.Create(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
