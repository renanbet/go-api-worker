package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/application/usecase"
)

type Handler struct {
	CreateOrderUC usecase.CreateOrder
	GetOrderUC    usecase.GetOrder
}

type createOrderRequest struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type createOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	res, err := h.CreateOrderUC.Execute(ctx, req.Product, req.Quantity)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, createOrderResponse{
		OrderID: res.OrderID,
		Status:  string(res.Status),
	})
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	o, err := h.GetOrderUC.Execute(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound{}) {
			writeJSONError(w, http.StatusNotFound, "order not found")
			return
		}
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, o)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
