package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dominicgisler/imap-spam-cleaner/database"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
)

func RunSummary(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var maxAge int
	if r.URL.Query().Has("maxage") {
		maxAge, err = strconv.Atoi(r.URL.Query().Get("maxage"))
		if err != nil {
			http.Error(w, "could not parse maxage", http.StatusBadRequest)
			logx.Errorf("could not parse maxage: %s", err)
			return
		}
	}

	inbox := r.URL.Query().Get("inbox")
	summaries, err := database.ListRunSummaries(inbox, maxAge)
	if err != nil {
		http.Error(w, "could not load run summaries", http.StatusInternalServerError)
		logx.Errorf("could not load run summaries: %s", err)
		return
	}

	var res []byte
	if inbox != "" && len(summaries) > 0 {
		res, err = json.Marshal(summaries[0])
	} else {
		res, err = json.Marshal(summaries)
	}

	if err != nil {
		http.Error(w, "could not load run summary", http.StatusInternalServerError)
		logx.Errorf("could not load run summary: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, "could not write response", http.StatusInternalServerError)
		logx.Errorf("could not write response: %s", err)
	}
}
