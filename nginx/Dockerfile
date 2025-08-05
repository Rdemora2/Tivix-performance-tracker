FROM node:22-alpine AS builder

WORKDIR /app
COPY ./front-end/package.json ./front-end/package-lock.json* ./
RUN npm ci
COPY ./front-end/ .
RUN npm run build

FROM nginx:1.27-alpine

RUN apk --no-cache add gettext

RUN rm /etc/nginx/conf.d/default.conf

RUN mkdir -p /etc/nginx/templates

COPY ./nginx/nginx.conf.template /etc/nginx/templates/nginx.conf.template

COPY ./nginx/entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/dist /usr/share/nginx/html

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]