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
	Title string `json:"title"`
	File  string `json:"file,omitempty"`
	Type  int    `json:"type"`
}

type CredentialItemData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TextItemData struct {
	Text string `json:"text"`
}

type FileItemData struct {
	File string `json:"file"`
}

type CreditCardItemData struct {
	Number string `json:"number"`
	Owner  string `json:"owner"`
	Period string `json:"period"`
	CVV    string `json:"CVV"`
}
