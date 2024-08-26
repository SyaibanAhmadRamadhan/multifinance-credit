package restapi_test

import (
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/product"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func Test_restApi_V1ProductsGet(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()

	mockProductService := product.NewMockService(mock)
	h := restapi.New(&service.Dependency{
		ProductService: mockProductService,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedParam := api.V1ProductsGetParams{
			Page:       1,
			PageSize:   10,
			MerchantId: nil,
			Ids:        nil,
		}

		expectedOutput := product.GetAllOutput{
			Pagination: pagination.PaginationOutput{
				Page:      1,
				PageSize:  10,
				PageCount: 1,
				TotalData: 1,
			},
			Items: []product.GetAllOutputItem{
				{
					ID:         rand.Int63(),
					MerchantID: 1,
					Name:       faker.Name(),
					Image:      faker.UUIDDigit(),
					Qty:        rand.Int63(),
					Price:      rand.Float64(),
				},
			},
		}

		params := url.Values{}
		params.Add("page", strconv.FormatInt(expectedParam.Page, 10))
		params.Add("page_size", strconv.FormatInt(expectedParam.PageSize, 10))

		r := chi.NewRouter()
		r.Get("/api/v1/product", h.V1ProductsGet)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/product?"+params.Encode(), nil)
		require.NoError(t, err)

		mockProductService.EXPECT().
			GetAll(req.Context(), product.GetAllInput{
				Pagination: pagination.PaginationInput{
					Page:     1,
					PageSize: 10,
				},
			}).
			Return(expectedOutput, nil)

		expectedRespon := api.V1ProductsGetResponseBody{
			Items: []api.V1ProductsGetResponseBodyItem{
				{
					Id:         expectedOutput.Items[0].ID,
					Image:      expectedOutput.Items[0].Image,
					MerchantId: expectedOutput.Items[0].MerchantID,
					Name:       expectedOutput.Items[0].Name,
					Price:      expectedOutput.Items[0].Price,
					Qty:        expectedOutput.Items[0].Qty,
				},
			},
			Pagination: api.PaginationResponse{
				Page:      1,
				PageCount: 1,
				PageSize:  10,
				TotalData: 1,
			},
		}

		rr := httptest.NewRecorder()
		h.V1ProductsGet(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.V1ProductsGetResponseBody
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, expectedRespon, resp)
	})
}
