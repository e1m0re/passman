package model

const (
	DatumTypeCredentials int = iota
	DatumTypeText
	DatumTypeFile
	DatumTypeCreditCard
)

type DatumInfo struct {
	UserID   int    `db:"user"`
	Metadata string `db:"metadata"`
	File     string `db:"file"`
	Checksum string `db:"checksum"`
}

type DatumItem struct {
	ID       int    `db:"id"`
	UserID   int    `db:"user"`
	Metadata string `db:"metadata"`
	File     string `db:"file"`
	Checksum string `db:"checksum"`
}

type DatumItemsList = []*DatumItem

type DatumMetadata struct {
	Title    string `json:"title"`
	FileName string `json:"fileName,omitempty"`
	Type     int    `json:"type"`
}

type CredentialItemData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
