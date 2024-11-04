package interface_repo

type IRelationRepo interface {
	GetFollowerAndFollowingCountofUser(userId *string) (*uint, *uint, *error)
	InitiateFollowRelationship(userId, userBId *string) (bool, error)
	InitiateUnFollowRelationship(userId, userBId *string) error
	GetFollowersIdsOfUser(userId *string) (*[]uint64, error)
	GetFollowingsIdsOfUser(userId *string) (*[]uint64, error)
	UserAFollowingUserBorNot(userId, userBId *string) (bool, error)
}
