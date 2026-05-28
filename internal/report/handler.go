package report

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jorgemorais/employee-api/internal/employee"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) StartReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := employee.ReportFilter{
		Type:     q.Get("type"),
		Position: q.Get("position"),
	}
	if v := q.Get("min_salary"); v != "" {
		filter.MinSalary, _ = strconv.ParseFloat(v, 64)
	}
	if v := q.Get("max_salary"); v != "" {
		filter.MaxSalary, _ = strconv.ParseFloat(v, 64)
	}

	reportID, err := h.svc.StartReport(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"report_id": reportID})
}

func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	reportID := mux.Vars(r)["report_id"]

	rep, ready, err := h.svc.GetReport(reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ready {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "processing"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rep)
}
