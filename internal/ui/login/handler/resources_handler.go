package handler

import (
	"io"
	"net/http"

	"github.com/caos/zitadel/internal/domain"
)

type dynamicResourceData struct {
	OrgID         string `schema:"orgId"`
	DefaultPolicy bool   `schema:"default-policy"`
}

func (l *Login) handleResources(staticDir http.FileSystem) http.Handler {
	return http.FileServer(staticDir)
}

func (l *Login) handleDynamicResources(w http.ResponseWriter, r *http.Request) {
	data := new(dynamicResourceData)
	err := l.getParseData(r, data)

	bucketName := domain.IAMID
	if data.OrgID != "" && !data.DefaultPolicy {
		bucketName = data.OrgID
	}
	reader, _, _ := l.staticStorage.GetObject(r.Context(), bucketName, domain.CssPath+"/"+domain.CssVariablesFileName)
	if err != nil {

	}
	bytes, _ := io.ReadAll(reader)
	w.Write(bytes)
}
