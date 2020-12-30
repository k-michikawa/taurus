package injectors

import (
	"taurus/domains/repositories"
	"taurus/infrastructures"
	infraRepo "taurus/infrastructures/repositories"
	"taurus/interfaces"
	"taurus/usecases"
)

// 抽象と実装の依存解決
func InjectUserRepository(db *infrastructures.Database) repositories.UserRepository {
	return infraRepo.NewUserRepository(*db)
}

func InjectUserUsecase(db *infrastructures.Database) usecases.UserUsecase {
	return usecases.NewUserUsecase(InjectUserRepository(db))
}

func InjectUserService(db *infrastructures.Database) interfaces.UserService {
	return interfaces.NewService(InjectUserUsecase(db))
}
