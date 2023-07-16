FROM golang:1.20.5-alpine AS build

COPY . /app
WORKDIR /app
RUN go build -o vigo360-server .

FROM alpine:3.18.2 AS sass

RUN apk add --no-cache nodejs npm
RUN npm install -g sass

COPY ./styles /app/styles

WORKDIR /app/styles

RUN sass --no-source-map --style compressed admin.scss admin.css
RUN sass --no-source-map --style compressed main.scss main.css

FROM alpine:3.18.2 AS runtime

RUN apk add --no-cache nginx supervisor
COPY ./docker/nginx.conf /etc/nginx/nginx.conf
COPY ./docker/supervisord.conf /etc/supervisord.conf

WORKDIR /app

COPY ./static /app/assets

COPY --from=build /app/vigo360-server /app/executable
COPY --from=sass /app/styles/admin.css /app/styles/main.css /app/assets/

ARG PORT=6000
ENV PORT=${PORT}

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
