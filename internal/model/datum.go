package model

type DatumTypeID = int

const (
	CredentialItem DatumTypeID = iota
	TextItem
	BinaryItem
	CreditCardItem
)

type DatumInfo struct {
	TypeID   DatumTypeID `db:"type"`
	UserID   int         `db:"user"`
	File     string      `db:"file"`
	Checksum string      `db:"checksum"`
}

type DatumItem struct {
	ID       int         `db:"id"`
	TypeID   DatumTypeID `db:"type"`
	UserID   int         `db:"user"`
	File     string      `db:"file"`
	Checksum string      `db:"checksum"`
}

type DatumItemsList = []*DatumItem
