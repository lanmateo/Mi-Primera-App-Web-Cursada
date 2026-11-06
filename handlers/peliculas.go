package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	db "catalogo-peliculas/db/sqlc"
	"catalogo-peliculas/views"
)

// PeliculaHandler agrupa los handlers HTTP que renderizan las vistas con templ.
type PeliculaHandler struct {
	queries *db.Queries
}

func NewPeliculaHandler(q *db.Queries) *PeliculaHandler {
	return &PeliculaHandler{queries: q}
}

// Index atiende GET /
// Obtiene todas las películas con sqlc y renderiza la página completa con templ.
func (h *PeliculaHandler) Index(w http.ResponseWriter, r *http.Request) {
	peliculas, err := h.queries.ListPeliculas(r.Context())
	if err != nil {
		log.Printf("error al listar películas: %v", err)
		http.Error(w, "Error al obtener las películas", http.StatusInternalServerError)
		return
	}

	if err := views.IndexPage("Catálogo de Películas", peliculas).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar la página: %v", err)
		http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
	}
}

// Create atiende POST /peliculas
// Lee el formulario, guarda la película con sqlc y redirige a / (patrón Post/Redirect/Get).
func (h *PeliculaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	if titulo == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	duracion, err := strconv.ParseInt(r.FormValue("duracion"), 10, 32)
	if err != nil {
		http.Error(w, "La duración debe ser un número entero", http.StatusBadRequest)
		return
	}

	anio, err := strconv.ParseInt(r.FormValue("anio_lanzamiento"), 10, 32)
	if err != nil {
		http.Error(w, "El año de lanzamiento debe ser un número entero", http.StatusBadRequest)
		return
	}

	puntaje, err := strconv.ParseFloat(r.FormValue("puntaje"), 64)
	if err != nil {
		http.Error(w, "El puntaje debe ser un número", http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreatePelicula(r.Context(), db.CreatePeliculaParams{
		Titulo:          titulo,
		Descripcion:     strings.TrimSpace(r.FormValue("descripcion")),
		Duracion:        int32(duracion),
		Actores:         strings.TrimSpace(r.FormValue("actores")),
		Genero:          strings.TrimSpace(r.FormValue("genero")),
		AnioLanzamiento: int32(anio),
		Puntaje:         puntaje,
	})
	if err != nil {
		log.Printf("error al crear película: %v", err)
		http.Error(w, "Error al guardar la película", http.StatusInternalServerError)
		return
	}

	// 303 See Other: el navegador hace un GET / y ve la lista actualizada
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete atiende POST /peliculas/{id}/eliminar
// Borra la película y redirige a / (también Post/Redirect/Get, sin JavaScript).
func (h *PeliculaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de película inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeletePelicula(r.Context(), int32(id)); err != nil {
		log.Printf("error al eliminar película: %v", err)
		http.Error(w, "Error al eliminar la película", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
