package responsemodel_chat

type NewGroupInfo struct {
	GroupName    string `json:"GroupName" `
	GroupMembers string `json:"GroupMembers"`
}

type AddNewMembersToGroupResponse struct {
	GroupID      string `json:"GroupID" `
	GroupMembers string `json:"GroupMembers" `
}

type RemoveMemberFromGroup struct {
	GroupID  string
	MemberID string
}
