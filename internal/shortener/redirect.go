package shortener

import (
	"net/http"
	"strings"
	"url-shortener/internal/db"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Extract the short code from the URL path path (e.g., /r/xf83g -> xf83g)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	code := pathParts[2]

	// Find the matching URL profile
	urlRecord, err := h.DB.GetURLByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	// Fire-and-forget record click metrics asynchronously to avoid blocking the user redirection pipeline
	go func(urlID int32, ip, ua string) {
		_ = h.DB.RecordClick(r.Context(), db.RecordClickParams{
			UrlID:     urlID,
			IpAddress: ip,
			UserAgent: ua,
		})
	}(urlRecord.ID, r.RemoteAddr, r.UserAgent())

	// Send standard HTTP 302 Found redirect
	http.Redirect(w, r, urlRecord.OriginalUrl, http.StatusFound)
}
