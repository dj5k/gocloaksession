package gocloaksession

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	gocloakHostname     = "http://localhost:8998"
	gocloakRealm        = "gocloak"
	gocloakClientID     = "gocloak"
	gocloakClientSecret = "gocloak-secret"
)

func initializeSession(t testing.TB) *goCloakSession {
	session, err := NewSession(gocloakClientID, gocloakClientSecret, gocloakRealm, gocloakHostname)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return session.(*goCloakSession)
}

func Test_GetKeycloakAuthToken_Authentication(t *testing.T) {
	session := initializeSession(t)

	token, err := session.GetKeycloakAuthToken()

	assert.NoError(t, err, "Login failed")
	assert.NotZero(t, token.AccessToken, "Token is not set")
}

func Test_GetKeycloakAuthToken_StillValid(t *testing.T) {
	session := initializeSession(t)

	_ = session.authenticate()

	require.NotNil(t, session.token, "Token is not set")
	require.NotZero(t, session.token.AccessToken, "Token is not set")
	require.NotZero(t, session.token.RefreshToken, "Token is not set")

	oldToken := session.token.AccessToken

	token, err := session.GetKeycloakAuthToken()

	assert.NoError(t, err, "refreshToken failed")
	assert.Equal(t, oldToken, token.AccessToken, "New AccessToken given, but expecting the old is still valid")
}

func Test_GetKeycloakAuthToken_Refresh(t *testing.T) {
	session := initializeSession(t)

	_ = session.authenticate()

	require.NotNil(t, session.token, "Token is not set")
	require.NotZero(t, session.token.AccessToken, "Token is not set")
	require.NotZero(t, session.token.RefreshToken, "Token is not set")

	oldToken := session.token.AccessToken
	session.token.AccessToken = ""

	token, err := session.GetKeycloakAuthToken()

	assert.NoError(t, err, "refreshToken failed")
	assert.NotEqual(t, oldToken, token.AccessToken, "No new AccessToken given")
}

func Test_refreshToken(t *testing.T) {
	session := initializeSession(t)

	_ = session.authenticate()

	require.NotNil(t, session.token, "Token is not set")
	require.NotZero(t, session.token.AccessToken, "Token is not set")
	require.NotZero(t, session.token.RefreshToken, "Token is not set")

	oldToken := session.token.AccessToken
	err := session.refreshToken()

	assert.NoError(t, err, "refreshToken failed")
	assert.NotEqual(t, oldToken, session.token.AccessToken, "No new AccessToken given")
}

func Test_authenticate(t *testing.T) {
	session := initializeSession(t)

	err := session.authenticate()

	assert.NoError(t, err, "authenticate failed")
	assert.NotZero(t, session.token.AccessToken, "Token is not set")
}
