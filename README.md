# Gogen Using Redis PubSub or Kafka

In this repo we are demonstrating on how to use the Redis Pubsub or Kafka communication between application using the gogen framework that also apply the Clean Architecture

## Gogen Framework
For the Gogen Framework Structure, you can refer to here link

> https://github.com/mirzaakhena/gogen

## Application Architecture

The application consist of two parts
1. Client : Has a restapi interface to invoke the Redis PubSub or Kafka publisher
2. Server : Has a Redis PubSub or Kafka server to consume the request, process it and then return back to client (Redis PubSub or Kafka)

![gogen pubsub architecture](https://github.com/mirzaakhena/gogen_pubsub/blob/main/gogen_pubsub_architecture.png)

## Folder structure
```text
gogen_pubsub
├── application
│  ├── app_client.go
│  └── app_server.go
├── domain_demo
│  ├── controller
│  │  ├── kafkasubscriber
│  │  ├── redissubscriber
│  │  └── restapi
│  ├── gateway
│  │  ├── simpleprint
│  │  ├── kafkapublisher
│  │  └── redispublisher
│  └── usecase
│      ├── onmessagereceived
│      └── runmessagesend
├── main.go
└── shared
    └── model
       └── payload
          └── data.go
```

## How to run the application

1. After you git clone it, make sure to run the `go mod tidy` to download the dependency
2. Run the server application by `go run main.go server`
3. Run the client application by `go run main.go client`
4. invoke this api with curl, postman or use the file `http_runmessagesend.http` under `domain_demo/controller/restapi`

    ```
    POST http://localhost:8000/api/v1/runmessagesend
    {
      "message": "hello" 
    }
    ```
    See the terminal from server side, it will print out the message. Means that the message has been sent by client and received by server
    ```
    >>> hello
    ```

## How to switch technology between gRPC and GraphQL

For the server you may comment / uncomment this part (`application/app_server.go`)
```
//primaryDriver := redissubscriber.NewController(log, cfg)
primaryDriver := kafkasubscriber.NewController(log, cfg)
```

For the client you may comment / uncomment this part (`application/app_client.go`)
```
//datasource := redispublisher.NewGateway(log, appData, cfg)
datasource := kafkapublisher.NewGateway(log, appData, cfg)
```
