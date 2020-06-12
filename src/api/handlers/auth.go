package handlers

import (
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"nomsu-api/tezos"
)

type AuthInput struct {
	Msg    string `json:"msg" form:"msg"`
	Sig    string `json:"sig" form:"sig"`
	PubKey string `json:"pubKey" form:"pubKey"`
}

type AuthResult struct {
	JWT string `json:"token"`
}

// @description Authenticate a user with a signed tezos message
// @router /auth [post]
// @success 200 {object} AuthResult "JWT token"
// @param input body AuthInput true "Signed message"
// @produce json
func Auth(c echo.Context) error {
	input := new(AuthInput)
	if err := c.Bind(input); err != nil {
		return err
	}

	// validate signature
	msg, _ := hex.DecodeString(input.Msg)
	sig := tezosprotocol.Signature(input.Sig)
	pubKey := tezosprotocol.PublicKey(input.PubKey)
	pk, _ := pubKey.CryptoPublicKey()
	err := tezos.VerifyOperation(msg, sig, pk)

	if err != nil {
		return echo.ErrUnauthorized
	}

	// validate timestamp
	msgTime := msgToTime(input.Msg)
	delta := time.Since(msgTime)
	allowed, _ := time.ParseDuration("5m")

	// time sync issues?
	if delta < 0 {
		return echo.ErrUnauthorized
	}

	// signed auth message expired
	if delta > allowed {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["tzid"], _ = tezosprotocol.NewContractIDFromPublicKey(pubKey)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // @todo expires at end of this game window

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return err
	}

	r := &AuthResult{JWT: t}
	return c.JSON(http.StatusOK, r)
}

func msgToTime(msg string) time.Time {
	mBytes, _ := hex.DecodeString(msg)
	mInt := binary.LittleEndian.Uint64(mBytes)
	return time.Unix(int64(mInt/1000), 0)
}
