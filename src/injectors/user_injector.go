package injectors

import (
	"taurus/domains/repositories"
	"taurus/infrastructures"
	infraRepo "taurus/infrastructures/repositories"
	"taurus/interfaces"
	"taurus/use_cases"
)

// 抽象と実装の依存解決
func InjectUserRepository(db *infrastructures.Database) repositories.UserRepository {
	return infraRepo.NewUserRepository(*db)
}

func InjectUserUsecase(db *infrastructures.Database) use_cases.UserUsecase {
	return use_cases.NewUserUsecase(InjectUserRepository(db))
}

func InjectUserService(db *infrastructures.Database) interfaces.UserService {
	return interfaces.NewService(InjectUserUsecase(db))
}
