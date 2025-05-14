package dto

type SuccessDto struct {
	Success bool `json:"success" default:"true"`
}

func CreateSuccessDto() *SuccessDto {
	dto := new(SuccessDto)
	dto.Success = true
	return dto
}
