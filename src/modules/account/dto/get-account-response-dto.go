package accountModuleDto

import (
	"go-gin-test-job/src/database/entities"
)

type GetAccountResponseDto struct {
	Offset int          `json:"offset"`
	Count  int          `json:"count"`
	Total  int64        `json:"total"`
	List   []AccountDto `json:"list"`
}

func CreateGetAccountResponseDto(offset int, count int, total int64, accounts []*entities.Account) GetAccountResponseDto {
	var dto GetAccountResponseDto
	dto.Offset = offset
	dto.Count = count
	dto.Total = total
	dto.List = make([]AccountDto, 0)
	for _, account := range accounts {
		dto.List = append(dto.List, CreateAccountDto(account))
	}
	return dto
}
