package interface_usecase

import (
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
)

type IAdminUsecase interface {
	AdminLogin(*requestmodels.AdminLoginData) (*responsemodels.AdminLoginres, error)
	GetAllUsers(string, string) (*[]responsemodels.UserAdminResponse, error)
	BlcokUser(string) error
	UnblockUser(string) error
}
