package app

import (
	"encoding/json"
	"fmt"
	"github.com/e1m0re/passman/internal/model"
	"log/slog"
)

func (a *app) updateItemsListView() {
	a.itemsListView.Clear()
	idx := 0
	for _, item := range a.store.itemsList {
		metadata := &model.DatumMetadata{}
		// Todo обработать ошибки
		err := json.Unmarshal([]byte(item.Metadata), metadata)
		if err != nil {
			s := fmt.Sprintf("%s", err.Error())
			slog.Info(s)
			continue
		}
		a.itemsListView.AddItem(metadata.Title, "", rune(49+idx), nil)
		idx++
	}
}
