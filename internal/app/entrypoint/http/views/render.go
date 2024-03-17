package views

import (
	"bytes"
	"encoding/json"
	"github.com/vaberof/vk-internship-task/pkg/http/protocols/apiv1"
	"net/http"
)

type httpStatus int

func RenderJSON(w http.ResponseWriter, status httpStatus, payload *apiv1.Response) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	w.Write(buf.Bytes())
}
