package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Gokul-B12/money-txn/db/mock"
	db "github.com/Gokul-B12/money-txn/db/sqlc"
	"github.com/Gokul-B12/money-txn/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {

	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)

	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)

}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("Matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := RandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				mockStore.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }

				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				//require.Equal(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidUname",
			body: gin.H{
				"username":  "dfa(8)",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }

				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//require.Equal(t, recorder.Body, user)
			},
		},
		{
			name: "ShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "affs",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }

				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//require.Equal(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "dafdajkdsf",
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }

				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//require.Equal(t, recorder.Body, user)
			},
		},
		{
			name: "DuplicateUname",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }

				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				//require.Equal(t, recorder.Body, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//creating ctrl to create mockstore

			ctrl := gomock.NewController(t)
			mockStore := mockdb.NewMockStore(ctrl)

			server := newTestServer(t, mockStore)

			//after creating mockSore and server, we will create stub

			tc.buildStubs(mockStore)

			url := "/users"

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}

}

func RandomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	return db.User{
		Username:       util.RandomString(5),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(5),
		Email:          util.RandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, response *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(response)
	require.NoError(t, err)

	var gotUser db.User

	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
	require.Equal(t, user.CreatedAt, gotUser.CreatedAt)
	require.Empty(t, gotUser.HashedPassword)
	require.NotEmpty(t, user.HashedPassword)

}
