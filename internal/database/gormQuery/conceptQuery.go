package gormQuery

import (
	"fmt"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllConceptsQuery(concepts *[]models.Concept, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error {
	return gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("display ILIKE ? AND code ILIKE ?", "%"+meaning+"%", code+"%").Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&concepts).Error; err != nil {
			return err
		} else if len(*concepts) == 0 {
			*concepts = []models.Concept{}
			return nil
		}
		return nil
	})
}
