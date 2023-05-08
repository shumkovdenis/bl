import { ConnectRouter } from "@bufbuild/connect";
import { IntegrationService } from "@shumkovdenis/protobuf-schema/lib/integration/v1/api_connect";

export default (router: ConnectRouter) =>
  router.service(IntegrationService, {
    async getBalance(req, ctx) {
      console.log("req:traceparent", ctx.requestHeader.get("traceparent"));
      console.log("req:grpc-trace-bin", ctx.requestHeader.get("grpc-trace-bin"));
      console.log("res:traceparent", ctx.responseHeader.get("traceparent"));
      console.log("res:grpc-trace-bin", ctx.responseHeader.get("grpc-trace-bin"));
      return { balance: BigInt(1010) };
    }
  });
