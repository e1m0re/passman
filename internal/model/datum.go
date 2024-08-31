package model

type DatumTypeID = int

const (
	CredentialItem DatumTypeID = iota
	TextItem
	BinaryItem
	CreditCardItem
)

type DatumInfo struct {
	UserID   int         `db:"user"`
	TypeID   DatumTypeID `db:"type"`
	File     string      `db:"file"`
	Checksum string      `db:"checksum"`
}

type DatumItem struct {
	ID       int         `db:"id"`
	UserID   int         `db:"user"`
	TypeID   DatumTypeID `db:"type"`
	File     string      `db:"file"`
	Checksum string      `db:"checksum"`
}
