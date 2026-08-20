package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"crud/internal/models"
	"crud/internal/repository"
)

type UsuarioHandler struct {
	repo repository.UsuarioRepository
}

func NewUsuarioHandler(repo repository.UsuarioRepository) *UsuarioHandler {
	return &UsuarioHandler{repo: repo}
}

func (h *UsuarioHandler) HandleUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		usuarios, err := h.repo.ObtenerTodos()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(usuarios)

	case http.MethodPost:
		var u models.Usuario
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		creado, err := h.repo.Crear(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(creado)
	}
}

func (h *UsuarioHandler) HandleUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/usuarios/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var u models.Usuario
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.ID = id
		if err := h.repo.Actualizar(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		if err := h.repo.Eliminar(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
