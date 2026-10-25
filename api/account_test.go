package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_db "github.com/chp0304/simplebank/db/mock"
	db "github.com/chp0304/simplebank/db/sqlc"
	"github.com/chp0304/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {

	testCases := []struct {
		name         string
		acount       db.Account
		accountID    int64
		buildStubs   func(store *mock_db.MockStore, account db.Account)
		checkReponse func(t *testing.T, recoder *httptest.ResponseRecorder, account db.Account)
	}{
		{

			name:   "ok",
			acount: randomAccount(),
			buildStubs: func(store *mock_db.MockStore, account db.Account) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkReponse: func(t *testing.T, recoder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, account)
			},
		},
		{
			name:   "not found",
			acount: randomAccount(),
			buildStubs: func(store *mock_db.MockStore, account db.Account) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkReponse: func(t *testing.T, recoder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, db.Account{})
			},
		},
		{
			name:   "internal error",
			acount: randomAccount(),
			buildStubs: func(store *mock_db.MockStore, account db.Account) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, fmt.Errorf("statusInternalServerError"))
			},
			checkReponse: func(t *testing.T, recoder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, db.Account{})
			},
		},
		{
			name: "bad request",
			acount: db.Account{
				ID:       0,
				Owner:    util.RandomOwner(),
				Currency: util.RandomCurrency(),
				Balance:  util.RandomInt(1, 100),
			},

			buildStubs: func(store *mock_db.MockStore, account db.Account) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0).Return(db.Account{}, fmt.Errorf("statusInternalServerError"))
			},
			checkReponse: func(t *testing.T, recoder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, db.Account{})
			},
		},
	}

	for _, tcase := range testCases {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		store := mock_db.NewMockStore(ctrl)

		server := NewServer(store)

		tcase.buildStubs(store, tcase.acount)

		recoder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d", tcase.acount.ID)

		request, err := http.NewRequest(http.MethodGet, url, nil)

		require.NoError(t, err)

		server.router.ServeHTTP(recoder, request)

		tcase.checkReponse(t, recoder, tcase.acount)
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomInt(1, 100),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
