FROM node:22.14.0-alpine AS build

ARG BACKEND_HOST

ENV VUE_APP_BACKEND_API_HOST=$BACKEND_HOST

WORKDIR /app
 
# Копируем package.json и package-lock.json (или yarn.lock)
COPY package*.json ./
RUN npm config set strict-ssl false && npm install
 
# Копируем исходный код
COPY . .
 
# Запускаем сборку приложения (например, Vue CLI создаёт папку dist)
RUN npm run build
 
# Этап 2. Запуск в продакшене через Nginx
FROM nginx:stable-alpine
# Копируем собранные файлы в папку, обслуживаемую Nginx
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
 
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
