package interfaces

import "github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"

type IAdminRepo interface {
	GetPassword(string) (string, error)
	AllUsers(limit, offset int) (*[]responsemodels.UserAdminResponse, error)
	BlockUser(string) error
	UnblockUser(string) error
}
