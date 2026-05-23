package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dominicgisler/imap-spam-cleaner/database"
)

type runSummaryResponse struct {
	Inbox     string                `json:"inbox,omitempty"`
	Summaries []database.RunSummary `json:"summaries"`
}

func RunSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inbox := r.URL.Query().Get("inbox")
	summaries, err := database.ListRunSummaries(inbox)
	if err != nil {
		http.Error(w, "could not load run summaries", http.StatusInternalServerError)
		return
	}

	var res []byte
	if inbox != "" && len(summaries) > 0 {
		res, err = json.Marshal(summaries[0])
	} else {
		res, err = json.Marshal(runSummaryResponse{
			Inbox:     inbox,
			Summaries: summaries,
		})
	}

	if err != nil {
		http.Error(w, "could not load run summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, "could not write response", http.StatusInternalServerError)
	}
}
