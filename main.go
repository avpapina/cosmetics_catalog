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

	// Раздача статических файлов из папки photos
	http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("./photos"))))

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleCatalogRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/catalog/")
	parts := strings.Split(path, "/")

	if len(parts) > 0 && (parts[0] == "favicon.ico" || parts[0] == "images" || parts[0] == "assets") {
		http.NotFound(w, r)
		return
	}

	switch {
	case len(parts) == 0 || parts[0] == "": // /catalog/
		handleMainPage(w, r)

	case len(parts) == 1 && parts[0] == "sales": // /catalog/sales
		handleSaleProducts(w, r)

	case len(parts) == 1 && parts[0] == "brands": // /catalog/brands
		handleBrands(w, r)

	case len(parts) == 1: // /catalog/{category}
		handleCatalogSubcategory(w, r, parts[0])

	case len(parts) == 2 && parts[0] == "sales": // /catalog/sales/{product}
		handleProduct(w, r, parts[0], "", parts[1])

	case len(parts) == 2 && parts[0] == "brands": // /catalog/brands/{brand}
		handleBrandProducts(w, r, parts[1])

	case len(parts) == 2: // /catalog/{category}/{subcategory}
		handleCategoryProducts(w, r, parts[0], parts[1])

	case len(parts) == 3 && parts[0] != "sales": // /catalog/{category}/{subcategory}/{product}
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

// Страница отображения брендов
func handleBrands(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/brands.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	var brands []models.Brand
	if err := database.DB.Find(&brands).Error; err != nil {
		http.Error(w, "Ошибка получения категорий", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, brands)
	if err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
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
	data := struct {
		Filter          string
		CategorySlug    string
		SubcategorySlug string
		MinPrice        float64
		MaxPrice        float64
		Products        []models.Product
	}{
		Filter:          filter,
		CategorySlug:    categorySlug,
		SubcategorySlug: subcategorySlug,
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
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
		IsOnSale    bool
		SalePrice   float64
	}{
		Name:        product.Name,
		Slug:        productSlug,
		Price:       product.Price,
		ImagePath:   product.ImagePath,
		Description: product.Description,
		IsOnSale:    product.IsOnSale,
		SalePrice:   product.SalePrice,
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

// Страница всех продуктов бренда
func handleBrandProducts(w http.ResponseWriter, r *http.Request, brandSlug string) {
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

	var brand models.Brand
	var products []models.Product
	var err error

	// Базовый запрос для получения бренда
	baseQuery := database.DB.Where("slug = ?", strings.ToLower(brandSlug))

	// В зависимости от фильтра применяем разные условия
	switch filter {
	case "no":
		// Без фильтрации (все продукты бренда)
		err = baseQuery.Preload("Products").First(&brand).Error
		products = brand.Products

	case "high":
		// Сортировка по возрастанию цены
		err = baseQuery.Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("products.price ASC")
		}).First(&brand).Error
		products = brand.Products

	case "low":
		// Сортировка по убыванию цены
		err = baseQuery.Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("products.price DESC")
		}).First(&brand).Error
		products = brand.Products

	case "range":
		// Фильтрация по ценовому диапазону
		err = baseQuery.Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("products.price BETWEEN ? AND ?", minPrice, maxPrice)
		}).First(&brand).Error
		products = brand.Products

	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Загружаем шаблон
	tmpl, err := template.ParseFiles("templates/products.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Подготавливаем данные для шаблона
	data := struct {
		Filter          string
		CategorySlug    string
		SubcategorySlug string
		MinPrice        float64
		MaxPrice        float64
		Products        []models.Product
	}{
		Filter:          filter,
		CategorySlug:    "brands",
		SubcategorySlug: brandSlug,
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
		Products:        products,
	}

	// Рендерим шаблон
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}

// Страница товаров со скидкой
func handleSaleProducts(w http.ResponseWriter, r *http.Request) {

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

	var products []models.Product
	var err error

	// Создаем базовый запрос
	dbQuery := database.DB.Model(&models.Product{}).Where("is_on_sale = true")

	// Применяем фильтры
	switch filter {
	case "no":
		// Без дополнительной фильтрации
		err = dbQuery.Find(&products).Error
	case "high":
		// Сортировка по возрастанию цены
		err = dbQuery.Order("price ASC").Find(&products).Error
	case "low":
		// Сортировка по убыванию цены
		err = dbQuery.Order("price DESC").Find(&products).Error
	case "range":
		// Фильтрация по ценовому диапазону
		err = dbQuery.
			Where("price BETWEEN ? AND ?", minPrice, maxPrice).
			Find(&products).Error
	}

	if err != nil {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Загружаем шаблон
	tmpl, err := template.ParseFiles("templates/products_sales.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Подготавливаем данные для шаблона
	data := struct {
		Filter   string
		MinPrice float64
		MaxPrice float64
		Products []models.Product
	}{
		Filter:   filter,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Products: products,
	}

	// Рендерим шаблон
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга: "+err.Error(), http.StatusInternalServerError)
	}
}
