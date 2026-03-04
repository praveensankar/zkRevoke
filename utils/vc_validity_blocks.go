package utils

func GetNumberOfBlocksVCisValid(validFrom int, validTo int) int {
	return (validTo - validFrom) / 86400
}
