package interfaces

import (
	"context"
	"database/sql"
	proto "taurus/infrastructures/proto"
	"taurus/usecases"
)

// proto.UserServiceServer これに従わないと意味がないので、要件を満たすため抽象を作る
type UserService interface {
	proto.UserServiceServer
}

// 実態はこっち
type userService struct {
	userUsecase usecases.UserUsecase
}

// returnの型を抽象として実態を返す
func NewService(userUsecase usecases.UserUsecase) UserService {
	return &userService{userUsecase}
}

// PostUser Service
func (s *userService) PostUser(c context.Context, r *proto.PostUserRequest) (*proto.PostUserResponse, error) {
	user, err := s.userUsecase.PostUser(r.Name, r.Email, r.Password)
	if err != nil {
		return nil, err
	}
	return &proto.PostUserResponse{
		User: &proto.User{
			Id:             user.ID.String(),
			Name:           user.Name,
			Email:          user.Email,
			CreatedAt:      user.CreatedAt.Unix(),
			UpdatedAtOneof: nil,
		},
	}, nil
}

// ListUser Service
func (s *userService) ListUser(c context.Context, r *proto.ListUserRequest) (*proto.ListUserResponse, error) {
	result, err := s.userUsecase.ListUser()
	if err != nil {
		return nil, err
	}
	users := make([]*proto.User, len(result))
	for _, user := range result {
		protoUser := &proto.User{
			Id:             user.ID.String(),
			Name:           user.Name,
			Email:          user.Email,
			CreatedAt:      user.CreatedAt.Unix(),
			UpdatedAtOneof: mapUpdatedAt(user.UpdatedAt),
		}
		users = append(users, protoUser)
	}
	return &proto.ListUserResponse{
		Users: users,
	}, nil
}

// ReadUser Service
func (s *userService) FindUser(c context.Context, r *proto.FindUserRequest) (*proto.FindUserResponse, error) {
	user, err := s.userUsecase.FindUser(r.Id)
	if err != nil {
		return nil, err
	}
	return &proto.FindUserResponse{
		User: &proto.User{
			Id:             user.ID.String(),
			Name:           user.Name,
			Email:          user.Email,
			CreatedAt:      user.CreatedAt.Unix(),
			UpdatedAtOneof: mapUpdatedAt(user.UpdatedAt),
		},
	}, nil
}

// UpdateUser Service
func (s *userService) UpdateUser(c context.Context, r *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	user, err := s.userUsecase.UpdateUser(r.Id, r.Name, r.Email, r.Password)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateUserResponse{
		User: &proto.User{
			Id:             user.ID.String(),
			Name:           user.Name,
			Email:          user.Email,
			CreatedAt:      user.CreatedAt.Unix(),
			UpdatedAtOneof: mapUpdatedAt(user.UpdatedAt),
		},
	}, nil
}

// DeleteUser Service
func (s *userService) DeleteUser(c context.Context, r *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if err := s.userUsecase.DeleteUser(r.Id); err != nil {
		return nil, err
	}
	return &proto.DeleteUserResponse{}, nil
}

func mapUpdatedAt(updatedAt sql.NullTime) *proto.User_UpdatedAt {
	if updatedAt.Valid {
		return &proto.User_UpdatedAt{
			UpdatedAt: updatedAt.Time.Unix(),
		}
	}
	return nil
}
