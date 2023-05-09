package main

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/gofiber/fiber/v2"
	integration "github.com/shumkovdenis/protobuf-schema/gen/integration/v1"
	integrationConnect "github.com/shumkovdenis/protobuf-schema/gen/integration/v1/integrationv1connect"
)

type InitInput struct {
	Token string `json:"token" validate:"required"`
}

type InitOutput struct {
	Balance int64 `json:"balance"`
}

type HTTPServer struct {
	integrationService integrationConnect.IntegrationServiceClient
}

func NewHTTPServer(cfg Config) error {
	integrationService := integrationConnect.NewIntegrationServiceClient(
		NewInsecureClient(),
		fmt.Sprintf("http://localhost:%d", cfg.Dapr.GRPCPort),
		connect.WithGRPC(),
		connect.WithInterceptors(
			NewTraceInterceptor(),
			NewLoggerInterceptor(),
			NewAppInterceptor(cfg.Integration.AppID),
		),
	)

	server := &HTTPServer{
		integrationService: integrationService,
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(
		NewTraceMiddleware(),
		NewLoggerMiddleware(),
	)
	app.Post("/init", server.Init)
	app.Post("/bet", server.Init)

	return app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (s *HTTPServer) Init(ctx *fiber.Ctx) error {
	var in InitInput

	if err := ctx.BodyParser(&in); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(&in); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ExtractValidateError(err))
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
