# syntax=docker/dockerfile:1.3
ARG DOCKER_IMAGE_PREFIX=
FROM ${DOCKER_IMAGE_PREFIX}node:17.3.0-alpine as build-stage

ARG BUILDKIT_INLINE_CACHE=1

WORKDIR /app
RUN npm install -g pnpm
COPY package*.json pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.pnpm-store \
    pnpm install
COPY . .
RUN --mount=type=cache,target=/root/.pnpm-store \
    pnpm run build

ARG DOCKER_IMAGE_PREFIX=
FROM ${DOCKER_IMAGE_PREFIX}nginx:1.21.4-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
