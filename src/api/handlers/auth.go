package handlers

import (
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"time"
	"strings"
	"strconv"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/gomodule/redigo/redis"

	"nomsu-api/tezos"
)

type AuthInput struct {
	Msg    	string `json:"msg" form:"msg"`
	Sig    	string `json:"sig" form:"sig"`
	PubKey 	string `json:"pubKey" form:"pubKey"`
	Address string `json:"address" form:"address"`
}

type AuthResult struct {
	JWT string `json:"token"`
}

var (
	conn redis.Conn
	err error
	reply interface{}
)

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

	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	// Handle player list (defines turn order)
	currentDay := time.Now().Format("2006-01-02")
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	playersListKey := "players:" + currentDay
	playerAddress := input.Address
	
	// Load existing logged in players
	players, err := redis.Strings(conn.Do("LRANGE", playersListKey, 0, -1))
	if err != nil {
		return err
	}

	// Find or create player Redis entry
	if len(players) == 0 {
		// No players exist, create entry
		CreatePlayerEntry(t, playerAddress, timestamp)
	} else {
		// Check if player has already logged
		// If so, their turn order will not be updated
		playerExists := PlayerExists(players, playerAddress)
		if !playerExists {
			// Player turn order does not exist, create entry
			CreatePlayerEntry(t, playerAddress, timestamp)
		}
	}

	r := &AuthResult{JWT: t}
	return c.JSON(http.StatusOK, r)
}

func msgToTime(msg string) time.Time {
	mBytes, _ := hex.DecodeString(msg)
	mInt := binary.LittleEndian.Uint64(mBytes)
	return time.Unix(int64(mInt/1000), 0)
}

func CreatePlayerEntry(jwt string, playerAddress string, timestamp string) {
	// Redis init
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}

	defer conn.Close()

	currentDay := time.Now().Format("2006-01-02")
	playersListKey := "players:" + currentDay

	var player struct {
		JWT  		string `redis:"JWT"`
		tzid 		string `redis:"tzid"`
		timestamp 	string `redis:"timestamp"`
	}

	player.JWT = jwt
	player.tzid = playerAddress
	player.timestamp = timestamp
	
	if _, err := conn.Do("LPUSH", playersListKey, player); err != nil {
		return
	}
}

func PlayerExists(slice []string, entry string) bool {
	for _, s := range slice {
		if strings.Contains(s, entry) {
			return true
		}
	}
    return false
}