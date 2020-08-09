FROM nginx:1.19-alpine
COPY dist/swagger-core-api/swagger.json /usr/share/nginx/html/docs/
COPY docs/index.html /usr/share/nginx/html/docs/
