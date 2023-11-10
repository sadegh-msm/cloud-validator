# cloud-validator

## API Service

### Overview

The API Service is responsible for receiving user information, including two images of the user. It stores the images in an S3 bucket, writes the user information to a database, and sends the information to a RabbitMQ queue for further processing.

### Setup

1. Clone the repository: `git clone https://github.com/sadegh-msm/cloud-validator`
2. Navigate to the `api-service` directory: `cd api-service`

### Configuration

Ensure that the necessary configurations are set in the `configs/config.yaml` file. This includes MongoDB connection details, S3 bucket credentials, and RabbitMQ configurations.

### Running the API Service

To run the API Service, execute the following commands:

```bash
cd api-service
go run main.go
```

The API Service will be accessible at `http://localhost:8080`.

### Dependencies

- ["github.com/aws/aws-sdk-go/aws"](https://"github.com/aws/aws-sdk-go/aws"): AWS sdk for Go.
- [github.com/labstack/echo](https://github.com/labstack/echo): Http library for Go.
- [github.com/rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go): rabbitMQ library for Go.
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus): Logging library for Go.

## Validation Service

### Overview

The Validation Service processes user information received from the API Service. It uses the information to send images to a face detection service, which then identifies face similarities. If faces are similar, an email alert is sent to the user.

### Setup

1. Clone the repository: `git clone https://github.com/sadegh-msm/cloud-validator`
2. Navigate to the `validation-service` directory: `cd validation-service`

### Configuration

Ensure that the necessary configurations are set in the `configs/config.yaml` file. This includes MongoDB connection details, S3 bucket credentials, and RabbitMQ configurations.

### Running the Validation Service

To run the Validation Service, execute the following commands:

```bash
cd validation-service
go run main.go
```

The Validation Service continuously validates user information and sends alerts when faces are similar.

### Dependencies

- ["github.com/aws/aws-sdk-go/aws"](https://"github.com/aws/aws-sdk-go/aws"): AWS sdk for Go.
- [github.com/labstack/echo](https://github.com/labstack/echo): Http library for Go.
- [github.com/rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go): rabbitMQ library for Go.
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus): Logging library for Go.

### Notes

- Both services use the same configuration file located at `configs/config.yaml`.
- It is essential to set up and configure the external services (S3, MongoDB, RabbitMQ) before running the services.

Feel free to reach out for any issues or improvements!