package impl

import (
  . "../../../auth-service"
  "github.com/common-go/auth"
  "time"
)

type UserInfoServiceImpl struct {
  UserService UserInfoService
}

func NewUserInfoServiceImpl() *UserInfoServiceImpl {
  return &UserInfoServiceImpl{}
}

func (s *UserInfoServiceImpl) GetUserInfo(info auth.AuthInfo) (*auth.UserInfo, error) {
  user := auth.UserInfo{}
  user.UserName = info.UserName
  now := time.Now()
  dur := now.Add(- time.Hour)
  user.AccessTimeTo = &now
  user.AccessTimeFrom = &dur
  return &user, nil
}

func (s *UserInfoServiceImpl) PassAuthentication(user auth.UserInfo) (bool, error) {
  return true, nil
}
func (s *UserInfoServiceImpl) HandleWrongPassword(user auth.UserInfo) (bool, error) {
  return true, nil
}
