FROM nginx:stable-alpine as base
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

FROM node:16-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM base as final
COPY --from=builder /app/dist /usr/share/nginx/html
