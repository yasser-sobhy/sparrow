package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/yasser-sobhy/sparrow/core"
)

// TokenAuth allow users to login using JWT token
type TokenAuth struct {
	APIEndpoint string
	Token       string
}

func NewTokenAuth() TokenAuth {
	return TokenAuth{
		ApiEndpoint: viper.GetInt("token_auth.api_endpoint"),
		Token:       viper.GetInt("token_auth.token"),
	}
}

func (tokenAuth *TokenAuth) Login(ws *core.Conn, inputToken []byte) bool {

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
}
