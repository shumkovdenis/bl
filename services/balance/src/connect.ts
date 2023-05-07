import { ConnectRouter } from "@bufbuild/connect";
import { IntegrationService } from "@shumkovdenis/protobuf-schema/lib/integration/v1/api_connect";

export default (router: ConnectRouter) =>
  router.service(IntegrationService, {
    async getBalance() {
        return { balance: BigInt(1010) };
    }
  });
