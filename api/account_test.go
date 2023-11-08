package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string                        // 테스트 케이스 이름
		accountID     int64                         // 얻고자하는 계정 ID
		buildStubs    func(store *mockdb.MockStore) //테스트 사례의 목적에 맞는 스텁 생성
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, account)
			},
		}, {
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		}, {
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		}, {
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},

		// CAN ADD MORE CASE
	}

	for i := range testCases {

		// 테스트 케이스
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stub
			tc.buildStubs(store)

			// start server and send requests
			// 테스트를 위해 실제 서버를 열 필요는 없다
			// httptest 패키지를 사용해서 api 요청의 응답을 기록하면 된다.
			server := NewServer(store)
			recoder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(t, recoder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		statusCode int
	}{
		{
			name: "OK",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}
				//build stubs
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			statusCode: http.StatusOK,
		}, {
			name: "BadRequest",
			body: gin.H{
				"owner":    account.Owner,
				"currency": "KRW",
			},
			buildStubs: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			statusCode: http.StatusBadRequest,
		}, {
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}
				//build stubs
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			statusCode: http.StatusInternalServerError,
		},
	}

	for i := range testCases {
		// 테스트 케이스
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stub
			tc.buildStubs(store)

			server := NewServer(store)
			recoder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			require.Equal(t, tc.statusCode, recoder.Code)
		})
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name       string
		accountID  int64
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		statusCode int
	}{
		{
			name:      "OK",
			accountID: account.ID,
			body: gin.H{
				"balance": account.Balance + 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateAccountParams{
					ID:      account.ID,
					Balance: account.Balance + 100,
				}
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			body: gin.H{
				"balance": account.Balance + 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateAccountParams{
					ID:      account.ID,
					Balance: account.Balance + 100,
				}
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			statusCode: http.StatusInternalServerError,
		}, {
			name:      "BadRequestByID",
			accountID: 0,
			body: gin.H{
				"balance": account.Balance + 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name:      "BadRequestByBalance",
			accountID: account.ID,
			body: gin.H{
				"balance": 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		// 테스트 케이스
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stub
			tc.buildStubs(store)

			server := NewServer(store)
			recoder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			require.Equal(t, tc.statusCode, recoder.Code)
		})
	}
}

func TestDeleteAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name       string
		accountID  int64
		buildStubs func(store *mockdb.MockStore)
		statusCode int
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(nil)
			},
			statusCode: http.StatusOK,
		}, {
			name:      "InternalErrorInGET",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			statusCode: http.StatusInternalServerError,
		}, {
			name:      "InternalErrorInDelete",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			statusCode: http.StatusInternalServerError,
		}, {
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			statusCode: http.StatusBadRequest,
		}, {
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			statusCode: http.StatusNotFound,
		},
	}

	for i := range testCases {
		// 테스트 케이스
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stub
			tc.buildStubs(store)

			server := NewServer(store)
			recoder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			require.Equal(t, tc.statusCode, recoder.Code)
		})
	}
}

func TestListAccountAPI(t *testing.T) {
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount()
	}

	type Query struct {
		pageID   int32
		pageSize int32
	}

	testCases := []struct {
		name          string                        // 테스트 케이스 이름
		query         Query                         // 얻고자하는 계정 ID
		buildStubs    func(store *mockdb.MockStore) //테스트 사례의 목적에 맞는 스텁 생성
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchAccounts(t, recoder.Body, accounts)
			},
		}, {
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		}, {
			name: "BadRequest",
			query: Query{
				pageID:   0,
				pageSize: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stub
			tc.buildStubs(store)

			server := NewServer(store)
			recoder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(t, recoder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	// account : 만든 계정
	// gotAccount : GET으로 얻은 계정
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	// accounts : 만든 계정들
	// gotAccounts : GET으로 얻은 계정들
	require.Equal(t, accounts, gotAccounts)
}
