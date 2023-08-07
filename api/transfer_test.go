package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Gokul-B12/money-txn/db/mock"
	db "github.com/Gokul-B12/money-txn/db/sqlc"
	"github.com/Gokul-B12/money-txn/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {

	account1 := RandomAccount()
	account2 := RandomAccount()

	amount := int64(10)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mockdb.NewMockStore(ctrl)

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

	//start test server and send request
	server := NewServer(mockStore)
	recorder := httptest.NewRecorder()

	url := "/transfers"

	body := gin.H{
		"from_account_id": account1.ID,
		"to_account_id":   account2.ID,
		"amount":          amount,
		"currency":        util.INR,
	}

	data, err := json.Marshal(body)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

}
