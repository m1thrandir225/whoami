FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

RUN npm install -g pnpm

# Install dependencies
RUN pnpm install --frozen-lockfile

COPY . .

ARG VITE_BACKEND_URL
ARG VITE_BACKEND_API_VERSION

ENV NODE_ENV=production
ENV DOCKER_BUILD=true
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL
ENV VITE_BACKEND_API_VERSION=$VITE_BACKEND_API_VERSION


RUN pnpm build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# Copy SSL certificates
RUN mkdir -p /etc/nginx/ssl

RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
