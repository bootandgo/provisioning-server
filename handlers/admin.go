package handlers

import (
	"net/http"
	"provisioning-server/models"
	storepkg "provisioning-server/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateRootAdminIfNeeded(store storepkg.Store, rootUsername, rootPassword string) {
	// Check if root admin already exists
	_, err := store.GetRootAdmin()
	if err == nil {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rootPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash root password")
	}

	rootAdmin := &models.Admin{
		Username:  rootUsername,
		Password:  string(hashedPassword),
		IsRoot:    true,
		CreatedAt: time.Now(),
	}

	if err := store.CreateAdmin(rootAdmin); err != nil {
		panic("failed to create root admin")
	}

}

func CreateInviteHandler(store storepkg.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only root admin can create invitations
		currentAdmin := c.MustGet("admin").(*models.Admin)
		if !currentAdmin.IsRoot {
			c.JSON(http.StatusForbidden, gin.H{"error": "only root admin can create invitations"})
			return
		}

		token := uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)

		invite := &models.Invitation{
			Token:     token,
			CreatedBy: currentAdmin.Username,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
			Used:      false,
		}

		if err := store.CreateInvitation(invite); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"token":      token,
			"expires_at": expiresAt.Format(time.RFC3339),
		})
	}
}

func AdminRegisterHandler(store storepkg.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Token    string `json:"token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate invitation token
		_, err := store.GetInvitation(req.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": storepkg.ErrInvalidInvitation.Error()})
			return
		}

		// Check if username already exists
		_, err = store.FindAdminByUsername(req.Username)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		admin := &models.Admin{
			Username:  req.Username,
			Password:  string(hashedPassword),
			IsRoot:    false,
			CreatedAt: time.Now(),
		}

		if err := store.CreateAdmin(admin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin"})
			return
		}

		// Mark invitation as used
		err = store.MarkInvitationUsed(req.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark invitation as used"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "admin created successfully"})
	}
}

func AdminLoginHandler(store storepkg.Store, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		admin, err := store.FindAdminByUsername(req.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": admin.Username,
			"exp":      time.Now().Add(72 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
