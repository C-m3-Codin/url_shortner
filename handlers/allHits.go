package handlers

import "github.com/gin-gonic/gin"

func GetHits(c *gin.Context) {

	// userID := c.GetUint("userID") // Replace with the actual user's ID.
	// fmt.Println("Hit all hits", userID)

	// // Query the user and preload their associated ShortLinks.
	// var user models.User

	// if result := services.DB.Preload("ShortLinks").First(&user, userID); result.Error != nil {
	// 	if result.Error == gorm.ErrRecordNotFound {
	// 		fmt.Println("User not found")
	// 	} else {
	// 		panic("Failed to load user details")
	// 	}
	// } else {
	// 	fmt.Printf("User: %s\n", user.Name)
	// 	fmt.Println("	:")
	// 	for _, shortLink := range user.ShortLinks {
	// 		fmt.Printf("- URL: %s\n", shortLink.OriginalURL)
	// 	}
	// }

}
