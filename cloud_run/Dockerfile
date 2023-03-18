FROM alpine:latest

RUN apk add nginx curl
RUN sed -i "s/return 404/return 200 Working./" /etc/nginx/http.d/default.conf

CMD ["nginx", "-g", "daemon off;"]

