package jwt

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	email  string
	userId string
)

func ProccesToken(token string, jwtSign string) (*models.Claim, bool, string, error) {
	var claim models.Claim
	key := []byte(jwtSign)
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return &claim, false, string(""), errors.New("Invalid token format.")
	}
	token = strings.TrimSpace(splitToken[1])
	// Decode jwt
	tkn, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err == nil {
		// Check database
		_, isOk, _ := database.CheckUser(claim.Email)
		if isOk {
			return &claim, isOk, claim.Id.Hex(), nil
		}
	}
	if !tkn.Valid {
		return &claim, false, string(""), errors.New("Invalid token.")
	}
	return &claim, false, string(""), err // err is nil
}

func GenerateJWT(ctx context.Context, user models.User) (string, error) {
	jwtSign := ctx.Value(models.Key("jwtSign")).(string)
	key := []byte(jwtSign)
	payload := jwt.MapClaims{
		"email":           user.Email,
		"nombre":          user.Nombre,
		"appellidos":      user.Apellidos,
		"fechaNacimiento": user.FechaNacimiento,
		"biografia":       user.Biografia,
		"ubicacion":       user.Ubicacion,
		"sitioWeb":        user.SitioWeb,
		"_id":             user.Id.Hex(),
		"exp":             time.Now().Add(time.Hour + 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return tokenStr, err
	}
	return tokenStr, nil
}
