FROM golang:1.21.7 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o go-jira .

RUN chmod +x ./go-jira

EXPOSE 8080

ENTRYPOINT [ "./go-jira" ]