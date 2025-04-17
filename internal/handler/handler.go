package handler

import (
	"html/template"
	"net/http"

	"myOrder/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	templates *template.Template
	db        *database.DB
}

type PageData struct {
	Error string
	Email string
}

func New(db *database.DB) (*Handler, error) {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		templates: templates,
		db:        db,
	}, nil
}

// Middleware to check authentication
func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value != "authenticated" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle the root path and login form submission
	if r.URL.Path == "/" {
		h.Landing(w, r)
		return
	}

	// Handle login POST request
	if r.URL.Path == "/login" && r.Method == http.MethodPost {
		h.Login(w, r)
		return
	}

	// Handle register routes
	if r.URL.Path == "/register" {
		if r.Method == http.MethodPost {
			h.Register(w, r)
		} else if r.Method == http.MethodGet {
			h.ShowRegister(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle logout
	if r.URL.Path == "/logout" {
		h.Logout(w, r)
		return
	}

	// Handle index page with authentication check
	if r.URL.Path == "/index.html" {
		cookie, err := r.Cookie("session")
		if err != nil || cookie == nil || cookie.Value != "authenticated" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		h.Index(w, r)
		return
	}

	// Handle any other paths
	http.NotFound(w, r)
}

func (h *Handler) Landing(w http.ResponseWriter, r *http.Request) {
	// Check if user is already authenticated
	cookie, err := r.Cookie("session")
	if err == nil && cookie != nil && cookie.Value == "authenticated" {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}

	data := PageData{}
	err = h.templates.ExecuteTemplate(w, "landing.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	_, err := h.db.GetUserByEmail(r.Context(), email)
	if err == nil {
		// Email already exists
		data := PageData{
			Error: "Email already registered. Please use a different email or login.",
			Email: email,
		}
		h.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	// Create new user
	_, err = h.db.CreateUser(r.Context(), email, password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Redirect to landing page after successful registration
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.db.GetUserByEmail(r.Context(), email)
		if err != nil {
			h.Landing(w, r)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			data := PageData{
				Error: "Invalid email or password",
				Email: email,
			}
			h.templates.ExecuteTemplate(w, "landing.html", data)
			return
		}

		// Create session cookie
		cookie := http.Cookie{
			Name:     "session",
			Value:    "authenticated",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   24 * 60 * 60, // 24 hours
		}
		http.SetCookie(w, &cookie)

		// Redirect to index page
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}

	// For GET requests, show the landing page
	h.Landing(w, r)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	cookie, err := r.Cookie("session")
	if err != nil || cookie == nil || cookie.Value != "authenticated" {
		// If no valid session cookie, redirect to login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = h.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie by setting MaxAge to -1
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // This will delete the cookie
	}
	http.SetCookie(w, &cookie)

	// Redirect to the landing page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
