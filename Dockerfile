FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG PORT==default_value
ARG GIN_MODE=default_value
ARG ALLOWED_ORIGINS=default_value
ARG REDIS_URI=default_value
ARG DATA_DB_HOST=default_value
ARG DATA_DB_USER=default_value
ARG DATA_DB_PASSWORD=default_value
ARG DATA_DB_PORT=default_value
ARG YOUTUBE_API_KEY1=default_value
ARG YOUTUBE_API_KEY2=default_value
ARG YOUTUBE_API_KEY3=default_value

ENV PORT=$PORT
ENV GIN_MODE=$GIN_MODE
ENV ALLOWED_ORIGINS=$ALLOWED_ORIGINS
ENV REDIS_URI=$REDIS_URI
ENV DATA_DB_HOST=$DATA_DB_HOST
ENV DATA_DB_USER=$DATA_DB_USER
ENV DATA_DB_PASSWORD=$DATA_DB_PASSWORD
ENV DATA_DB_PORT=$DATA_DB_PORT
ENV YOUTUBE_API_KEY1=$YOUTUBE_API_KEY1
ENV YOUTUBE_API_KEY2=$YOUTUBE_API_KEY2
ENV YOUTUBE_API_KEY3=$YOUTUBE_API_KEY3

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]
