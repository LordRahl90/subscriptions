FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o subscriptions_api ./cmd/api
RUN go build -o subscriptions_seed ./cmd/seed


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/subscriptions_api subscriptions_api
COPY --from=build /app/subscriptions_seed subscriptions_seed


EXPOSE 8080