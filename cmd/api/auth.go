package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID        primitive.ObjectID `json:"_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Role      string             `json:"role"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	//claims["sub"] = fmt.Sprint(user.ID)
	claims["sub"] = user.ID.Hex()
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	claims["role"] = user.Role

	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["role"] = user.Role

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return TokenPairs
	return tokenPairs, nil
}

func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (*jwt.Token, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	//get auth header
	authHeader := r.Header.Get("Authorization")

	//sanity check
	if authHeader == "" {
		return nil, nil, errors.New("no auth header")
	}

	//split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return nil, nil, errors.New("invalid auth header")
	}

	//check if we have bearer
	if headerParts[0] != "Bearer" {
		return nil, nil, errors.New("invalid auth header")
	}

	//extract token sting
	tokenStr := headerParts[1]

	//declare empty claims
	claims := &Claims{}

	//parse the token with claims
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, nil, errors.New("token is expired")
		}
		return nil, nil, err
	}

	if claims.Issuer != j.Issuer {
		return nil, nil, errors.New("invalid token issuer")
	}

	// Return the parsed token and claims
	return token, claims, nil
}

func (j *Auth) GetUserIDFromHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	w.Header().Add("Vary", "Authorization")

	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Sanity check for the Authorization header
	if authHeader == "" {
		return "", errors.New("no auth header")
	}

	// Split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid auth header")
	}

	// Check if we have a Bearer token
	if headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	// Extract the token string
	tokenStr := headerParts[1]

	// Declare empty claims
	claims := &Claims{}

	// Parse the token and extract claims
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	// Handle token parsing errors
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", errors.New("token is expired")
		}
		return "", err
	}

	// Verify the issuer
	if claims.Issuer != j.Issuer {
		return "", errors.New("invalid token issuer")
	}

	// Return the UserID (subclaim)
	return claims.Subject, nil
}
