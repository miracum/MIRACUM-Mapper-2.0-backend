package utilities

func GetOffset(page int, pageSize int) int {
	return (page - 1) * pageSize
}
