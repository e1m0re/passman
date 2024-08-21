package models

type UsersDataItemTypeID int

const (
	CredentialItem UsersDataItemTypeID = iota
	TextItem
	BinaryItem
	CreditCardItem
)

type UsersDataItemID int

type UsersDataItemInfo struct {
	TypeID   UsersDataItemTypeID `db:"type"`
	UserID   UserID              `db:"user"`
	File     []byte              `db:"file"`
	Checksum []byte              `db:"checksum"`
}

type UsersDataItem struct {
	ID       UsersDataItemID     `db:"id"`
	TypeID   UsersDataItemTypeID `db:"type"`
	UserID   UserID              `db:"user"`
	File     []byte              `db:"file"`
	Checksum []byte              `db:"checksum"`
}
