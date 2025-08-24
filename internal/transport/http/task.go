package http

import (
	"checklist-api-service/internal/dto"
	"encoding/json"
	"net/http"
)

func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errDTO := dto.NewErr(err.Error())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)

		return
	}

	
}
