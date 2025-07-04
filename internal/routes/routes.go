package routes

import (
	"juno-backend/configs"
	"juno-backend/internal/api" // ✅ Keep this - it works
	"juno-backend/internal/auth"
	"juno-backend/internal/middleware"

	// Remove: "juno-backend/internal/handlers"  // ❌ This doesn't exist

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *configs.Config) *gin.Engine {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "juno-backend",
			"message": "🚗 Juno Backend is running!",
		})
	})

	// Initialize OAuth before setting up routes
	auth.InitOAuth(cfg)

	// OAuth routes (no auth required)
	r.GET("/auth/google", auth.GoogleLogin(cfg))
	r.GET("/auth/google/callback", auth.GoogleCallback(cfg))

	// Protected routes (require JWT)
	protected := r.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// Auth endpoints
		protected.GET("/auth/me", auth.GetCurrentUser)
		protected.POST("/auth/logout", auth.Logout) // ✅ Add this logout route

		// API endpoints - Use the working api package functions
		protected.GET("/api/profile", api.GetProfile)
		protected.PUT("/api/profile", api.UpdateProfile)
		protected.GET("/api/friends", api.GetFriends)                    // ✅ Real friends list
		protected.POST("/api/friends", api.AddFriend)                    // ✅ Add friend by ID
		protected.GET("/api/friends/requests", api.GetFriendRequests)    // ✅ Pending requests
		protected.POST("/api/friends/username", api.AddFriendByUsername) // ✅ Add by username
		protected.GET("/api/users/search", api.SearchUsers)              // ✅ User search
		protected.GET("/api/rides", api.GetRides)
		protected.POST("/api/rides", api.CreateRide)
		protected.GET("/api/rides/nearby", api.GetNearbyRides)
		protected.GET("/api/rides/:id", api.GetRideDetails)
		protected.POST("/api/rides/:id/join", api.JoinRide)
		protected.DELETE("/api/rides/:id/leave", api.LeaveRide)
		protected.POST("/api/rides/:id/cancel", api.CancelRide)
	}

	return r
}
