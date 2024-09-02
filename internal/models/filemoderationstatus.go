package models

// FileModerationStatus represents the moderation status for files
type FileModerationStatus struct {
	ID             string `bson:"_id,omitempty"`
	Filename       string `bson:"filename"`
	FilePath       string `bson:"filepath"`
	ApprovalStatus string `bson:"approvalstatus"`
}
