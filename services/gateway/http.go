package main

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/shumkovdenis/bl/services/gateway/helpers"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type InitInput struct {
	Token string `json:"token" validate:"required"`
}

type InitOutput struct {
	Balance int64 `json:"balance"`
}

type HTTPServer struct {
	client             *req.Client
	integrationService integrationConnect.IntegrationServiceClient
}

func NewHTTPServer(cfg Config) error {
	client := req.C().
		SetBaseURL(fmt.Sprintf("http://localhost:%d", cfg.Dapr.HTTPPort)).
		WrapRoundTripFunc(
			helpers.NewClientLoggerMiddleware(),
			helpers.NewClientTraceMiddleware(),
			helpers.NewClientAppMiddleware(cfg.Integration.AppID),
		)

	integrationService := integrationConnect.NewIntegrationServiceClient(
		helpers.NewInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
		connect.WithInterceptors(
			helpers.NewAppInterceptor(cfg.Integration.AppID),
			helpers.NewTraceInterceptor(),
			helpers.NewLoggerInterceptor(),
		),
	)

	server := &HTTPServer{
		client:             client,
		integrationService: integrationService,
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(
		helpers.NewServerTraceMiddleware(),
		helpers.NewServerLoggerMiddleware(),
	)
	app.Post("/gateway", server.Gateway)
	app.Post("/init", server.Init)
	app.Post("/bet", server.Init)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (s *HTTPServer) Gateway(ctx *fiber.Ctx) error {
	data := fiber.Map{}

	res, err := s.client.R().
		SetContext(ctx.UserContext()).
		SetSuccessResult(&data).
		Post("/integration")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if !res.IsSuccessState() {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": res.String(),
		})
	}

	return ctx.JSON(&data)
}

func (s *HTTPServer) Init(ctx *fiber.Ctx) error {
	var in InitInput

	if err := ctx.BodyParser(&in); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := helpers.Validate(&in); len(err) != 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	req := connect.NewRequest(&integration.GetBalanceRequest{})

	res, err := s.integrationService.GetBalance(ctx.UserContext(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var out InitOutput
	out.Balance = res.Msg.Balance

	return ctx.JSON(&out)
}
