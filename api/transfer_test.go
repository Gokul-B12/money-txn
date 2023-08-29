package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Gokul-B12/money-txn/db/mock"
	db "github.com/Gokul-B12/money-txn/db/sqlc"
	"github.com/Gokul-B12/money-txn/token"
	"github.com/Gokul-B12/money-txn/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {

	user1, _ := RandomUser(t)
	user2, _ := RandomUser(t)
	user3, _ := RandomUser(t)

	account1 := RandomAccount(user1.Username)
	account2 := RandomAccount(user2.Username)
	account3 := RandomAccount(user3.Username)

	amount := int64(10)
	account1.Currency = util.INR
	account2.Currency = util.INR
	account3.Currency = util.USD

	testCases := []struct {
		name          string
		body          gin.H
		setAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				//Checking the existence of account1 in Accounts table
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        10,
				}

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "UnauthorisedUser",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "Unauthorised_user", time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				//Checking the existence of account1 in Accounts table
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        10,
				}

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "Unauthorised_user", time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
					//Return(account1, nil)

				//Checking the existence of account1 in Accounts table
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
					//Return(account2, nil)

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "fromAccountIDNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name: "toAccountIDNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},

		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(0)
					//Return(account3, nil)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).
					Return(account3, nil)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        "JUI",
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account3.ID,
				"amount":          -amount,
				"currency":        "INR",
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        "INR",
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "TransferTxnError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "INR",
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				//Checking the existence of account1 in Accounts table
				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				mockStore.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				// arg := db.TransferTxParams{
				// 	FromAccountID: account1.ID,
				// 	ToAccountID:   account2.ID,
				// 	Amount:        10,
				// }

				mockStore.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)

			tc.buildStubs(mockStore)

			//start test server and send request
			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			url := "/transfers"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}

}
