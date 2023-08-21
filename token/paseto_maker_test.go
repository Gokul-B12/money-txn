package token

import (
	"testing"
	"time"

	"github.com/Gokul-B12/money-txn/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {

	symmetricKey := util.RandomString(32)
	Maker, err := NewJWTMaker(symmetricKey)
	require.NoError(t, err)

	duration := time.Minute
	username := util.RandomOwner()

	token, err := Maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := Maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	//below vars are to check the returned payload data
	issuedAt := time.Now()
	exipiredAt := time.Now().Add(duration)
	require.Equal(t, payload.Username, username)
	require.NotZero(t, payload.ID, token)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, exipiredAt, time.Second)

}

//below test is to test with expired token

func TestExpiredPasetoToken(t *testing.T) {

	symmetricKey := util.RandomString(32)
	Maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	duration := -time.Minute
	username := util.RandomOwner()

	token, err := Maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := Maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}

// this case is to check token with none algorithm --- not possible in paseto
// func TestInvalidJWTTokenAlgoNone(t *testing.T) {
// 	symmetricKey := util.RandomString(32)
// 	Maker, err := NewJWTMaker(symmetricKey)
// 	require.NoError(t, err)

// 	duration := time.Minute
// 	username := util.RandomOwner()

// 	payload, err := NewPayload(username, duration)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, payload)

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
// 	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 	require.NoError(t, err)

// 	payload, err = Maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrInvalidToken.Error())
// 	require.Nil(t, payload)

// }
