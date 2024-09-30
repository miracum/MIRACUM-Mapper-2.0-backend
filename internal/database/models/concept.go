package models

type Concept struct {
	ID                  uint64 `gorm:"primarykey"`
	Code                string // `gorm:"index,gin"`
	Display             string // `gorm:"index,gin"`
	CodeSystemID        uint32 // `gorm:"index"`
	Elements            []Element
	CodeSystem          CodeSystem
	DisplaySearchVector string `gorm:"type:tsvector"` // Correctly map as tsvector for PostgreSQL
}

// func (c *Concept) BeforeSave(tx *gorm.DB) (err error) {
// 	tx.Exec("UPDATE concepts SET display_search_vector = to_tsvector(?) WHERE id = ?", c.Display, c.ID)
// 	return
// }

// Migration to add the tsvector column and index it
// func Migrate(db *gorm.DB) error {
// 	return db.Exec(`
//         ALTER TABLE concepts ADD COLUMN IF NOT EXISTS display_search_vector tsvector;
// 		GENERATED ALWAYS AS (to_tsvector('english', display)) STORED;
//         CREATE INDEX IF NOT EXISTS idx_display_search_vector ON concepts USING gin (display_search_vector);
//     `).Error
// }
