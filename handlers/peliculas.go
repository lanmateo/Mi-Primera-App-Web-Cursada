package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "catalogo-peliculas/db/sqlc"
)

// PeliculaDTO representa la estructura esperada en las peticiones JSON
type PeliculaDTO struct {
	Titulo          string  `json:"titulo"`
	Descripcion     string  `json:"descripcion"`
	Duracion        int32   `json:"duracion"`
	Actores         string  `json:"actores"`
	Genero          string  `json:"genero"`
	AnioLanzamiento int32   `json:"anio_lanzamiento"`
	Puntaje         float64 `json:"puntaje"`
}

type PeliculaHandler struct {
	queries *db.Queries
}

func NewPeliculaHandler(q *db.Queries) *PeliculaHandler {
	return &PeliculaHandler{queries: q}
}

// parseID extrae y valida el ID numérico de la URL (soporta Go 1.22 PathValue y parseo tradicional)
func parseID(r *http.Request) (int32, error) {
	if idVal := r.PathValue("id"); idVal != "" {
		id, err := strconv.Atoi(idVal)
		return int32(id), err
	}
	trimmed := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) >= 2 {
		id, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			return int32(id), nil
		}
	}
	return 0, fmt.Errorf("ID inválido")
}

// HandlePeliculas rutea peticiones a la colección (/peliculas)
func (h *PeliculaHandler) HandlePeliculas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlePeliculaByID rutea peticiones a recursos individuales (/peliculas/{id})
func (h *PeliculaHandler) HandlePeliculaByID(w http.ResponseWriter, r *http.Request) {
	// Si acceden a /peliculas/ sin ID, redirigir al handler de colección
	if r.URL.Path == "/peliculas/" || r.URL.Path == "/peliculas" {
		h.HandlePeliculas(w, r)
		return
	}

	id, err := parseID(r)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// POST /peliculas
func (h *PeliculaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input PeliculaDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validación: el título no puede estar vacío
	if strings.TrimSpace(input.Titulo) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	pelicula, err := h.queries.CreatePelicula(r.Context(), db.CreatePeliculaParams{
		Titulo:          input.Titulo,
		Descripcion:     input.Descripcion,
		Duracion:        input.Duracion,
		Actores:         input.Actores,
		Genero:          input.Genero,
		AnioLanzamiento: input.AnioLanzamiento,
		Puntaje:         input.Puntaje,
	})
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pelicula)
}

// GET /peliculas
func (h *PeliculaHandler) List(w http.ResponseWriter, r *http.Request) {
	peliculas, err := h.queries.ListPeliculas(r.Context())
	if err != nil {
		http.Error(w, "Error al listar películas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if peliculas == nil {
		peliculas = []db.Pelicula{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peliculas)
}

// GET /peliculas/{id}
func (h *PeliculaHandler) Get(w http.ResponseWriter, r *http.Request, id int32) {
	pelicula, err := h.queries.GetPelicula(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Película no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al buscar película: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pelicula)
}

// PUT /peliculas/{id}
func (h *PeliculaHandler) Update(w http.ResponseWriter, r *http.Request, id int32) {
	// Verificar existencia previa
	_, err := h.queries.GetPelicula(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Película no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al buscar película: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var input PeliculaDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validación
	if strings.TrimSpace(input.Titulo) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	err = h.queries.UpdatePelicula(r.Context(), db.UpdatePeliculaParams{
		ID:              id,
		Titulo:          input.Titulo,
		Descripcion:     input.Descripcion,
		Duracion:        input.Duracion,
		Actores:         input.Actores,
		Genero:          input.Genero,
		AnioLanzamiento: input.AnioLanzamiento,
		Puntaje:         input.Puntaje,
	})
	if err != nil {
		http.Error(w, "Error al actualizar película: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver el objeto actualizado
	actualizada, err := h.queries.GetPelicula(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al recuperar película actualizada: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actualizada)
}

// DELETE /peliculas/{id}
func (h *PeliculaHandler) Delete(w http.ResponseWriter, r *http.Request, id int32) {
	// Verificar existencia previa
	_, err := h.queries.GetPelicula(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Película no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al buscar película: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeletePelicula(r.Context(), id); err != nil {
		http.Error(w, "Error al eliminar película: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
