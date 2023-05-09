package main

import (
	"log"
)

/*
func newError(transactionID string) error {
	err := connect.NewError(
		connect.CodeInvalidArgument,
		errors.New("player id is required"),
	)

	rollbackInfo := &integration.RollbackInfo{
		TransactionId: transactionID,
	}

	if detail, detailErr := connect.NewErrorDetail(rollbackInfo); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func (s *Server) GetBalance(
	ctx context.Context,
	req *connect.Request[integration.GetBalanceRequest],
) (*connect.Response[integration.GetBalanceResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, err // automatically coded correctly
	}

	if req.Msg.PlayerId == "" {
		return nil, newError("123")
	}

	log.Println("req:traceparent", req.Header().Get("traceparent"))
	log.Println("req:tracestate", req.Header().Get("tracestate"))
	log.Println("req:grpc-trace-bin", req.Header().Get("grpc-trace-bin"))

	// client, err := dapr.NewClient()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// in := &dapr.InvokeBindingRequest{
	// 	Name:      s.walletBindingName,
	// 	Operation: "post",
	// 	Data:      []byte(""),
	// 	Metadata:  map[string]string{"path": "/6b9663d1-41a3-47f8-8e56-8e5c8678bcde"},
	// }

	// event, err := client.InvokeBinding(ctx, in)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, connect.NewError(
	// 		connect.CodeInvalidArgument,
	// 		err,
	// 	)
	// }

	// log.Println(event.Metadata)
	// log.Println(string(event.Data))
	// log.Println(event.Metadata["statusCode"])

	// data := BalanceData{}
	// if err := json.Unmarshal(event.Data, &data); err != nil {
	// 	return nil, err
	// }

	// traceparent, err := tracing.Parse(req.Header().Get("traceparent"))
	// if err != nil {
	// 	return nil, err
	// }

	// span, err := traceparent.NewSpan()
	// if err != nil {
	// 	return nil, err
	// }

	r := connect.NewRequest(&integration.GetBalanceRequest{PlayerId: "123"})
	r.Header().Set("dapr-app-id", "remote")
	r.Header().Set("traceparent", req.Header().Get("tracestate"))
	// r.Header().Set("traceparent", span.String())
	r.Header().Set("grpc-trace-bin", req.Header().Get("grpc-trace-bin"))

	log.Println("balance-req:traceparent", r.Header().Get("traceparent"))
	log.Println("balance-req:grpc-trace-bin", r.Header().Get("grpc-trace-bin"))

	t, err := s.client.GetBalance(ctx, r)
	if err != nil {
		log.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			log.Println(connectErr.Message())
			log.Println(connectErr.Details())
		}
		return nil, err
	}

	log.Println("balance-res:traceparent", t.Header().Get("traceparent"))
	log.Println("balance-res:grpc-trace-bin", t.Header().Get("grpc-trace-bin"))

	res := connect.NewResponse(&integration.GetBalanceResponse{
		Balance: t.Msg.Balance,
	})
	// res.Header().Set("traceparent", span.String())

	log.Println("res:traceparent", res.Header().Get("traceparent"))
	log.Println("res:grpc-trace-bin", res.Header().Get("grpc-trace-bin"))

	return res, nil
}
*/

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server started on port %d in %s mode", cfg.Port, cfg.Mode)

	if cfg.Mode == "grpc" {
		if err := NewGRPCServer(cfg); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := NewHTTPServer(cfg); err != nil {
			log.Fatal(err)
		}
	}
}
