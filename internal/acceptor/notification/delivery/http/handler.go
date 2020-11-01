package http

import (
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification"
	"github.com/KonstantinPronin/email-sending-service/pkg/constants"
	"github.com/KonstantinPronin/email-sending-service/pkg/middleware"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"math"
	"strconv"
)

type AcceptorHandler struct {
	logger    *zap.Logger
	server    *echo.Echo
	sanitizer *bluemonday.Policy
	usecase   notification.Usecase
}

func NewAcceptorHandler(
	server *echo.Echo,
	sanitizer *bluemonday.Policy,
	usecase notification.Usecase,
	logger *zap.Logger) {
	handler := AcceptorHandler{
		logger:    logger,
		server:    server,
		sanitizer: sanitizer,
		usecase:   usecase,
	}

	server.GET("/notifs/:id", handler.GetMessage, middleware.ParseErrors)
	server.GET("/notifs", handler.GetPage, middleware.ParseErrors)
	server.POST("/notifs", handler.AcceptMessage, middleware.ParseErrors)
}

func (handler *AcceptorHandler) GetMessage(ctx echo.Context) error {
	id := ctx.Param("id")

	notif, err := handler.usecase.Get(id)
	if err != nil {
		return err
	}

	if _, err := easyjson.MarshalToWriter(notif, ctx.Response().Writer); err != nil {
		handler.logger.Error(fmt.Sprintf("Response marshal error: %s", err.Error()))
		return err
	}

	return nil
}

func (handler *AcceptorHandler) GetPage(ctx echo.Context) error {
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(ctx.QueryParam("per_page"))
	if err != nil || perPage < 1 {
		perPage = constants.DefaultPageSize
	}

	list, total, err := handler.usecase.GetList(int64(page), int64(perPage))
	if err != nil {
		return err
	}

	handler.addPageHeader(page, perPage, int(total), ctx)

	if _, err := easyjson.MarshalToWriter(list, ctx.Response().Writer); err != nil {
		handler.logger.Error(fmt.Sprintf("Response marshal error: %s", err.Error()))
		return err
	}

	return nil
}

func (handler *AcceptorHandler) AcceptMessage(ctx echo.Context) error {
	notif := new(model.Notification)

	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, notif); err != nil {
		return model.NewInvalidArgument("wrong request body format")
	}

	notif.Sender = handler.sanitizer.Sanitize(notif.Sender)
	notif.Subject = handler.sanitizer.Sanitize(notif.Subject)
	notif.Message = handler.sanitizer.Sanitize(notif.Message)

	for i, val := range notif.To {
		notif.To[i] = handler.sanitizer.Sanitize(val)
	}

	id, err := handler.usecase.Accept(notif)
	if err != nil {
		return err
	}

	if _, err := easyjson.MarshalToWriter(model.Notification{ID: id}, ctx.Response().Writer); err != nil {
		handler.logger.Error(fmt.Sprintf("Response marshal error: %s", err.Error()))
		return err
	}

	return nil
}

func (handler *AcceptorHandler) addPageHeader(page, perPage, total int, ctx echo.Context) {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	ctx.Response().Header().Set(constants.Total, strconv.Itoa(total))
	ctx.Response().Header().Set(constants.TotalPages, strconv.Itoa(totalPages))
	ctx.Response().Header().Set(constants.PerPage, strconv.Itoa(perPage))
	ctx.Response().Header().Set(constants.Page, strconv.Itoa(page))

	if page+1 < totalPages {
		ctx.Response().Header().Set(constants.NextPage, strconv.Itoa(page+1))
	}

	if page-1 > 0 {
		ctx.Response().Header().Set(constants.PrevPage, strconv.Itoa(page-1))
	}
}
