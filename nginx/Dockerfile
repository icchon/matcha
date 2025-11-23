FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY web/index.html /usr/share/nginx/html/index.html

CMD ["nginx", "-g", "daemon off;"]
