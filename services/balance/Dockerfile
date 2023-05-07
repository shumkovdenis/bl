FROM node:18-alpine AS build

ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=$GITHUB_TOKEN
RUN echo //npm.pkg.github.com/:_authToken=$GITHUB_TOKEN >> ~/.npmrc

RUN npm install -g pnpm

WORKDIR /app
COPY .npmrc pnpm-lock.yaml ./
RUN pnpm fetch --prod

ADD . ./
RUN pnpm install --offline --prod

RUN pnpm run build

FROM gcr.io/distroless/nodejs:18
COPY --from=build /app/bundle.js /usr/src/app/
WORKDIR /usr/src/app
CMD ["bundle.js"]
