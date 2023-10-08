package article

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-bunrouter-gorm-example/infrastructure/httplib"
	logger "go-bunrouter-gorm-example/infrastructure/log"
	"go-bunrouter-gorm-example/infrastructure/validator"
	"go-bunrouter-gorm-example/module/primitive"
	"go-bunrouter-gorm-example/utils"

	"github.com/uptrace/bunrouter"
)

type Http struct {
	serviceArticle InterfaceService
}

func NewHttp(serviceHealth InterfaceService) InterfaceHttp {
	return &Http{
		serviceArticle: serviceHealth,
	}
}

type InterfaceHttp interface {
	GroupArticle(group *bunrouter.Group)
}

func (h *Http) GroupArticle(g *bunrouter.Group) {
	g.GET("", h.GetListArticle)
	g.GET("/:id", h.DetailArticle)
	g.POST("", h.CreateArticle)
}

func (h *Http) GetListArticle(w http.ResponseWriter, c bunrouter.Request) error {
	logCtx := fmt.Sprintf("handler.GetListArticle")
	ctx := context.Background()

	if h.serviceArticle == nil {
		err := errors.New("dependency service article to handler article on method GetListArticle is nil")
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceHealth")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}

	paginationQuery, err := httplib.GetPaginationFromCtx(c)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "httplib.GetPaginationFromCtx")
		return httplib.SetErrorResponse(w, http.StatusBadRequest, err.Error())
	}

	query := c.Request.URL.Query().Get("query")
	if query != "" {
		if !utils.IsValidSanitizeSQL(query) {
			err = errors.New(primitive.QueryIsSuspicious)
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "utils.IsValidSanitizeSQL")
			return httplib.SetErrorResponse(w, http.StatusBadRequest, primitive.QueryIsSuspicious)
		}
	}

	author := c.Request.URL.Query().Get("author")
	if author != "" {
		if !utils.IsValidSanitizeSQL(author) {
			err = errors.New(primitive.QueryIsSuspicious)
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "utils.IsValidSanitizeSQL")
			return httplib.SetErrorResponse(w, http.StatusBadRequest, primitive.QueryIsSuspicious)
		}
	}

	param := primitive.ParameterArticleHandler{
		Query:  query,
		Author: author,
	}

	data, count, err := h.serviceArticle.GetListArticle(ctx, param, paginationQuery)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceArticle.GetListArticle")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}

	return httplib.SetPaginationResponse(w,
		http.StatusOK,
		primitive.SuccessGetArticle,
		data,
		uint64(count),
		paginationQuery)

}

func (h *Http) CreateArticle(w http.ResponseWriter, c bunrouter.Request) error {
	logCtx := fmt.Sprintf("handler.CreateArticle")
	ctx := context.Background()

	if h.serviceArticle == nil {
		err := errors.New("dependency service article to handler article on method CreateArticle is nil")
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceHealth")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}

	var requestBody primitive.ArticleReq
	// Decode the request body into the Article struct.
	if err := json.NewDecoder(c.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return httplib.SetErrorResponse(w, http.StatusBadRequest, primitive.SomethingWrongWithTheBodyRequest)
	}

	errValidateStruct := validator.ValidateStructResponseSliceString(requestBody)
	if errValidateStruct != nil {
		logger.Error(ctx, logCtx, "validator.ValidateStructResponseSliceString got err : %v", errValidateStruct)
		return httplib.SetCustomResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil, errValidateStruct)
	}

	data, err := h.serviceArticle.RecordArticle(ctx, requestBody)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceArticle.GetListArticle")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}

	return httplib.SetSuccessResponse(w, http.StatusOK, primitive.SuccessCreateArticle, data)

}

func (h *Http) DetailArticle(w http.ResponseWriter, c bunrouter.Request) error {
	logCtx := fmt.Sprintf("handler.DetailArticle")
	ctx := context.Background()

	if h.serviceArticle == nil {
		err := errors.New("dependency service article to handler article on method DetailArticle is nil")
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceHealth")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}

	idParam := c.Param("id")
	if idParam == "" {
		err := errors.New(primitive.ParamIdIsZeroOrNullString)
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "c.Param")
		return httplib.SetErrorResponse(w, http.StatusBadRequest, primitive.ParamIdIsZeroOrNullString)
	}

	idInt64, err := strconv.Atoi(idParam)
	if err != nil || idInt64 == 0 {
		err := errors.New(primitive.ParamIdIsZeroOrNullString)
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "strconv.Atoi")
		return httplib.SetErrorResponse(w, http.StatusBadRequest, primitive.ParamIdIsZeroOrNullString)
	}

	data, err := h.serviceArticle.GetDetailArticle(ctx, int64(idInt64))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceArticle.GetDetailArticle")
			return httplib.SetErrorResponse(w, http.StatusNotFound, primitive.RecordArticleNotFound)
		}
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceArticle.GetDetailArticle")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}

	return httplib.SetSuccessResponse(w, http.StatusOK, primitive.SuccessCreateArticle, data)

}
