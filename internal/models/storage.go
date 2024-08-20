package models

type UsersDataItemTypeID int

const (
	CredentialItem UsersDataItemTypeID = iota
	TextItem
	BinaryItem
	CreditCardItem
)

type UsersDataID int

type UsersDataItem struct {
	ID       UsersDataID         `db:"id"`
	TypeID   UsersDataItemTypeID `db:"type"`
	UserID   UserID              `db:"user"`
	File     []byte              `db:"file"`
	Checksum []byte              `db:"checksum"`
}
