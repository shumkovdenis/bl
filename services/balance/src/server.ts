import { fastify } from "fastify";
import { fastifyConnectPlugin } from "@bufbuild/connect-fastify";
import routes from "./connect";
import { parseEnv, port } from "znv";

async function main() {  
  const { PORT } = parseEnv(process.env, {
    PORT: port().default(6000),
  });

  const server = fastify({http2: true});
  await server.register(fastifyConnectPlugin, {
    routes,
  });
  server.get("/", (_, reply) => {
    reply.type("text/plain");
    reply.send("Hello World!");
  });
  await server.listen({ host: "localhost", port: PORT });
  console.log("server is listening at", server.addresses());
}
// You can remove the main() wrapper if you set type: module in your package.json,
// and update your tsconfig.json with target: es2017 and module: es2022.
void main();
