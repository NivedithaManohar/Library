package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}
		tokenString := authHeader // Token without "Bearer" prefix

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract role from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			return
		}

		// Set role in Gin context
		c.Set("role", role)
		fmt.Println("RRRRRRRRRRRRRRRRRRRRR", role)

		// Continue processing the request
		c.Next()
	}
}

// type Claims struct {
// 	UserRole string `json:"user_role"`
// 	jwt.StandardClaims
// }

// func ExtractUserRoleFromToken(tokenString string) (string, error) {
// 	// Parse the token
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("SECRET")), nil // Change "your-secret-key" to your actual secret key
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	// Check if the token is valid
// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		return claims.UserRole, nil
// 	}

// 	return "", fmt.Errorf("invalid token")
// }

// // Middleware to extract user role from JWT token and set it in the Gin context
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Extract JWT token from Authorization header
// 		authHeader := c.GetHeader("Authorization")
// 		fmt.Println("Authorization header:", authHeader) // Debugging statement
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
// 			return
// 		}
// 		authHeaderParts := strings.Split(authHeader, " ")
// 		fmt.Println("Authorization header parts:", authHeaderParts) // Debugging statement
// 		if len(authHeaderParts) < 2 {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
// 			return
// 		}

// 		tokenString := authHeaderParts[1]
// 		fmt.Println("JWT token:", tokenString) // Debugging statement
// 		// Extract user role from JWT token
// 		userRole, err := ExtractUserRoleFromToken(tokenString)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Set user role in Gin context
// 		c.Set("userRole", userRole)

// 		fmt.Println("RRRRRRRRRRRRRRRRRRRRRRRRRRRr", userRole)

// 		// Continue processing the request
// 		c.Next()
// 	}
// }
