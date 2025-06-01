package main

import (
	"cosmetics_catalog/database"
	"cosmetics_catalog/models"
	"cosmetics_catalog/repositories"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var productRepo *repositories.ProductRepository

func main() {
	// Подключение к базе данных
	if err := database.Connect(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Посев тестовых данных
	if err := database.SeedTestData(); err != nil {
		log.Fatalf("Ошибка при посеве данных: %v", err)
	}

	// Инициализация репозитория продуктов
	productRepo = repositories.NewProductRepository(database.DB)

	// Настройка маршрутов
	http.HandleFunc("/", handleCatalogRoutes)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleCatalogRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/catalog/")
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 0 || parts[0] == "":
		handleMainPage(w, r)
	case len(parts) == 1 && parts[0] != "" && parts[0] != "sales": // /
		handleCatalogSubcategory(w, r, parts[0])
	case len(parts) == 1 && parts[0] != "" && parts[0] == "sales": // /catalog/sales/
		handleCategoryProducts(w, r, parts[0], "")
	case len(parts) == 2 && parts[0] != "sales": // /catalog/{category}
		handleCategoryProducts(w, r, parts[0], parts[1])
	case len(parts) == 2 && parts[0] == "sales" && parts[1] != "": // /catalog/sales/{product}
		handleProduct(w, r, parts[0], "", parts[1])
	case len(parts) == 3 && parts[1] != "" && parts[0] != "sales": // /catalog/{category}/{subcategory}/{product}
		handleProduct(w, r, parts[0], parts[1], parts[2])

	default:
		http.NotFound(w, r)
	}
}

// Главная страница с категориями
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/catalog.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		http.Error(w, "Ошибка получения категорий", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, categories)
	if err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}

// Страница подкатегории макияж уход
func handleCatalogSubcategory(w http.ResponseWriter, r *http.Request, slug string) {
	var current models.Category

	// Находим категорию и подгружаем подкатегории
	if err := database.DB.
		Preload("Subcategories").
		Where("slug = ?", slug).
		First(&current).Error; err != nil {

		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/subcategory.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Создаём структуру данных для шаблона
	data := struct {
		Name          string
		Slug          string
		Subcategories []models.Subcategory
	}{
		Name:          current.Name,
		Slug:          current.Slug,
		Subcategories: current.Subcategories,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}

// Страница всех продуктов подкатегории
func handleCategoryProducts(w http.ResponseWriter, r *http.Request, categorySlug, subcategorySlug string) {

	// Получаем параметры фильтрации
	query := r.URL.Query()
	filter := query.Get("filter")
	if filter == "" {
		filter = "no" // Значение по умолчанию
	}

	// Обработка ценового диапазона
	var minPrice, maxPrice float64
	if filter == "range" {
		minPrice, _ = strconv.ParseFloat(query.Get("min_price"), 64)
		maxPrice, _ = strconv.ParseFloat(query.Get("max_price"), 64)
	}

	// Отображение продуктов без фильтра
	var subcat models.Subcategory
	if filter == "no" {
		if err := database.DB.
			Preload("Products").
			Joins("JOIN categories ON categories.id = subcategories.category_id").
			Where("subcategories.slug = ? AND categories.slug = ?",
				strings.ToLower(subcategorySlug),
				strings.ToLower(categorySlug)).
			First(&subcat).Error; err != nil {

			http.NotFound(w, r)
			return
		}
	}

	// Сортировка продуктов по возрастанию цены
	if filter == "high" {
		if err := database.DB.
			Preload("Products", func(db *gorm.DB) *gorm.DB {
				return db.Order("products.price ASC")
			}).
			Joins("JOIN categories ON categories.id = subcategories.category_id").
			Where("subcategories.slug = ? AND categories.slug = ?",
				strings.ToLower(subcategorySlug),
				strings.ToLower(categorySlug)).
			First(&subcat).Error; err != nil {

			http.NotFound(w, r)
			return
		}
	}

	// Сортировка продуктов по убыванию цены
	if filter == "low" {
		if err := database.DB.
			Preload("Products", func(db *gorm.DB) *gorm.DB {
				return db.Order("products.price DESC")
			}).
			Joins("JOIN categories ON categories.id = subcategories.category_id").
			Where("subcategories.slug = ? AND categories.slug = ?",
				strings.ToLower(subcategorySlug),
				strings.ToLower(categorySlug)).
			First(&subcat).Error; err != nil {

			http.NotFound(w, r)
			return
		}
	}

	// Отображение продуктов в промежутке от minPrice до maxPrice
	if filter == "range" {
		if err := database.DB.
			Preload("Products", func(db *gorm.DB) *gorm.DB {
				return db.
					Where("products.price BETWEEN ? AND ?", minPrice, maxPrice)
			}).
			Joins("JOIN categories ON categories.id = subcategories.category_id").
			Where("subcategories.slug = ? AND categories.slug = ?",
				strings.ToLower(subcategorySlug),
				strings.ToLower(categorySlug)).
			First(&subcat).Error; err != nil {

			http.NotFound(w, r)
			return
		}
	}

	// Загружаем шаблон
	tmpl, err := template.ParseFiles("templates/products.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Подготавливаем данные для шаблона
	// Подготавливаем данные для шаблона
	data := struct {
		Filter          string
		CategorySlug    string
		SubcategorySlug string
		MinPrice        float64 // Добавляем
		MaxPrice        float64 // Добавляем
		Products        []models.Product
	}{
		Filter:          filter,
		CategorySlug:    categorySlug,
		SubcategorySlug: subcategorySlug,
		MinPrice:        minPrice, // Добавляем
		MaxPrice:        maxPrice, // Добавляем
		Products:        subcat.Products,
	}

	// Рендерим шаблон
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}

// Страница конкретного продукта
func handleProduct(w http.ResponseWriter, r *http.Request, category, subcategory, productSlug string) {

	var product models.Product
	if err := database.DB.Where("slug = ?", productSlug).First(&product).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	// Создаем структуру данных для шаблона
	data := struct {
		Name        string
		Slug        string
		Price       float64
		ImagePath   string
		Description string
	}{
		Name:        product.Name,
		Slug:        productSlug,
		Price:       product.Price,
		ImagePath:   product.ImagePath,
		Description: product.Description,
	}

	tmpl, err := template.ParseFiles("templates/product.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return

	}

	// Рендерим шаблон
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}
