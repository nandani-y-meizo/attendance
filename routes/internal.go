package routes

// import (
// 	"net/http"

// 	"shared/pkgs/jwtmanager"

// 	"github.com/gin-gonic/gin"
// )

// // JwtStatus returns status information about JWT public key availability
// func JwtStatus(c *gin.Context) {
// 	init := jwtmanager.IsInitialized()
// 	vault := jwtmanager.IsVaultBacked()
// 	fp, _ := jwtmanager.PublicKeyFingerprint()
// 	c.JSON(http.StatusOK, gin.H{
// 		"initialized":  init,
// 		"vault_backed": vault,
// 		"fingerprint":  fp,
// 	})
// }
