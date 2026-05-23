FROM node:22-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

FROM deps AS dev

COPY vite.admin.config.js ./
COPY web/ ./web/

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]

FROM deps AS build

COPY vite.admin.config.js ./
COPY web/ ./web/
RUN npm run build

FROM nginx:1.27-alpine AS production

COPY docker/admin.nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist/admin/ /usr/share/nginx/html/

EXPOSE 80
