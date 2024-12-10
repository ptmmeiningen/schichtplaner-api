package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shift-planner/api/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("Login Decode Error: %v", err)
		return
	}

	var admin models.Admin
	if result := h.db.Where("username = ?", loginReq.Username).First(&admin); result.Error != nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Ungültige Anmeldedaten")
		log.Printf("Login Error: Benutzer %s nicht gefunden", loginReq.Username)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginReq.Password)); err != nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Ungültige Anmeldedaten")
		log.Printf("Login Error: Falsches Passwort für Benutzer %s", loginReq.Username)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id":      admin.ID,
		"username":      admin.Username,
		"is_super_user": admin.IsSuperUser,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler bei der Token-Generierung")
		log.Printf("Login Token Generation Error: %v", err)
		return
	}

	SendSuccessResponse(w, "Login erfolgreich", LoginResponse{Token: tokenString})
	log.Printf("Login: Benutzer %s erfolgreich angemeldet", admin.Username)
}

func (h *AdminHandler) GetAdmins(w http.ResponseWriter, r *http.Request) {
	var admins []models.Admin
	result := h.db.Find(&admins)

	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Abrufen der Administratoren")
		log.Printf("GetAdmins DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Administratoren erfolgreich abgerufen", admins)
	log.Printf("GetAdmins: %d Administratoren abgerufen", len(admins))
}
