package accountModuleDto

import (
	"go-gin-test-job/src/database/entities"
)

func CreatePostCreateAccountResponseDto(account *entities.Account) AccountDto {
	return CreateAccountDto(account)
}
