package main

import (
	"cosmetics_catalog/database"
	"cosmetics_catalog/models"
	"cosmetics_catalog/repositories"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	// Получаем подкатегорию с продуктами
	var subcat models.Subcategory
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

	// Загружаем шаблон
	tmpl, err := template.ParseFiles("templates/products.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Подготавливаем данные для шаблона
	data := struct {
		CategorySlug    string
		SubcategorySlug string
		Products        []models.Product
	}{
		CategorySlug:    categorySlug,
		SubcategorySlug: subcategorySlug,
		Products:        subcat.Products, // Используем загруженные продукты
	}

	// Рендерим шаблон
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}

// Страница конкретного продукта
func handleProduct(w http.ResponseWriter, r *http.Request, category, subcategory, productSlug string) {
	tmpl, err := template.ParseFiles("templates/product.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	product, err := productRepo.GetBySlug(productSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = tmpl.Execute(w, product)
	if err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}
