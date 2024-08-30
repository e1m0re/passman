package model

type DatumTypeID = int

const (
	CredentialItem DatumTypeID = iota
	TextItem
	BinaryItem
	CreditCardItem
)

type DatumId string
type UserID int

type DatumInfo struct {
	UserID   UserID      `db:"user"`
	TypeID   DatumTypeID `db:"type"`
	File     []byte      `db:"file"`
	Checksum []byte      `db:"checksum"`
}

type DatumItem struct {
	ID       DatumId     `db:"id"`
	UserID   UserID      `db:"user"`
	TypeID   DatumTypeID `db:"type"`
	File     []byte      `db:"file"`
	Checksum []byte      `db:"checksum"`
}
