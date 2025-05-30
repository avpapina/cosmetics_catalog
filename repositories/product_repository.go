package repositories

import (
	"cosmetics_catalog/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository создает новый экземпляр репозитория
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create добавляет новый продукт
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID возвращает продукт по ID
func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Brand").
		Preload("Subcategory.Category").
		First(&product, id).
		Error
	return &product, err
}

// GetOnSale возвращает товары со скидкой
func (r *ProductRepository) GetProductsOnSale() ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Where("is_on_sale = ?", true).
		Preload("Brand").
		Preload("Subcategory.Category").
		Find(&products).
		Error
	return products, err
}

func (r *ProductRepository) GetByCategory(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Joins("JOIN subcategories ON products.subcategory_id = subcategories.id").
		Where("subcategories.category_id = ?", categoryID).
		Preload("Brand").
		Preload("Subcategory.Category").
		Find(&products).
		Error
	return products, err
}

// Update обновляет продукт с проверкой существования
func (r *ProductRepository) Update(product *models.Product) error {
	// Проверяем, существует ли продукт
	if err := r.db.First(&models.Product{}, product.ID).Error; err != nil {
		return err // Продукт не найден
	}
	return r.db.Save(product).Error
}

// Delete удаляет продукт с проверкой
func (r *ProductRepository) Delete(id uint) error {
	// Проверка существования перед удалением
	if err := r.db.First(&models.Product{}, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Product{}, id).Error
}

// SearchByName поиск по названию с пагинацией
func (r *ProductRepository) SearchByName(query string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Where("name LIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Preload("Brand").
		Preload("Subcategory.Category").
		Find(&products).
		Error
	return products, err
}
