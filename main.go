package main

import (
	"html/template"
	"net/http"
	"strings"
)

func main() {

	// Обработчик главной страницы (каталога)
	http.HandleFunc("/", handleCatalogRoutes)

	//Обработчик перехода из каталога
	//http.HandleFunc("/catalog/", handleCatalogSubcategory)

	// Запускаем сервер
	http.ListenAndServe(":8080", nil)
}

// Обработчик главной страницы
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/catalog.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, catalog)
	if err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}

// Обработчик подкатегорий каталога
func handleCatalogSubcategory(w http.ResponseWriter, r *http.Request, slug string) {

	var current Category
	for i := range catalog {
		if catalog[i].Slug == slug {
			current = catalog[i]
			break
		}
	}

	if current.Name == "" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/subcategory.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, current)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}

// Обработчик продуктов подкатегории
func handleCategoryProducts(w http.ResponseWriter, r *http.Request, category string, subcategory string) {

	tmpl, err := template.ParseFiles("templates/products.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	data := struct {
		CategorySlug    string
		SubcategorySlug string
		Products        []Product
	}{
		CategorySlug:    category,
		SubcategorySlug: subcategory,
		Products:        products,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}

}

// Обработчик продукта
func handleProduct(w http.ResponseWriter, r *http.Request, category string, subcategory string, product string) {

	var foundProduct *Product
	for _, p := range products {
		if p.Name == product {
			foundProduct = &p
			break
		}
	}

	tmpl, err := template.ParseFiles("templates/product.html")

	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, foundProduct)

	if err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}

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
