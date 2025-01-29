# go-super-farmer

## Single Service

### Environment Variables

Create a .env file in the root directory of the project and add the following variables:

```bash
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_PORT=5432
DB_TIMEZONE=Asia/Jakarta
JWT_SECRET_KEY=secret
RABBITMQ_HOST=localhost
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_PORT=5672
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis
SMTP_HOST=localhost
SMTP_PORT=1025
EMAIL_FROM=noreply@example.com
EMAIL_PASSWORD=password
REPORT_PORT=8081
SERVER_PORT=8080
```

### Build

#### With Docker

```bash
docker build -t go-super-farmer .
```

#### Without Docker

```bash
go build -o go-super-farmer .
```

### Run

```bash
./go-super-farmer
```

## With Docker Compose to run All Services

1. Clone all Super Farmer Services in one directory

2. Create a docker compose.yml file
Directory Structure
```bash
.
├── docker compose.yml
├── go-super-farmer-api
├── go-super-farmer-report-service
├── go-super-farmer-mail-service
└── predict-model
```

3. Copy the docker compose.yml file
```bash
#for production

services:
  rabbitmq:
    image: "rabbitmq:4.0.5-management-alpine"
    container_name: rabbitmq
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:15672/api/overview"]
      interval: 10s
      timeout: 5s
      retries: 5
    ports:
      - "5672:5672"
      - "15672:15672"
    networks:
      - app-network
    volumes:
      - "rabbitmq_data:/data"

  postgres:
    image: postgres:latest
    container_name: postgres
    environment:
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=${DB_NAME}
    ports:
      - "5432:5432"
    networks:
      - app-network
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:latest
    container_name: redis
    ports:
      - "6379:6379"
    networks:
      - app-network
    restart: unless-stopped

  api:
    build:
      context: api
      dockerfile: Dockerfile
      target: runner
    depends_on:
      - rabbitmq
      - postgres
      - redis
    ports:
      - "${SERVER_PORT:-8080}:${SERVER_PORT:-8080}"
    networks:
      - app-network
    restart: unless-stopped
    environment:
      - DB_HOST=${DB_HOST}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_PORT=${DB_PORT}
      - DB_TIMEZONE=${DB_TIMEZONE}
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - RABBITMQ_HOST=${RABBITMQ_HOST}
      - RABBITMQ_USER=${RABBITMQ_USER}
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - RABBITMQ_PORT=${RABBITMQ_PORT}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - EMAIL_FROM=${EMAIL_FROM}
      - EMAIL_PASSWORD=${EMAIL_PASSWORD}
      - REPORT_PORT=${REPORT_PORT}
      - SERVER_PORT=${SERVER_PORT}
      - CASBIN_MODEL_PATH=/configs/model.conf
      - CASBIN_POLICY_PATH=/configs/policy.csv
    volumes:
      - ./api/pkg/auth/casbin/model.conf:/configs/model.conf
      - ./api/pkg/auth/casbin/policy.csv:/configs/policy.csv

  mail:
    build:
      context: mail
      dockerfile: Dockerfile
      target: runner
    depends_on:
      - rabbitmq
    networks:
      - app-network
    restart: unless-stopped
    environment:
      - RABBITMQ_HOST=${RABBITMQ_HOST}
      - RABBITMQ_USER=${RABBITMQ_USER}
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - RABBITMQ_PORT=${RABBITMQ_PORT}

  report:
    build:
      context: report
      dockerfile: Dockerfile
      target: runner
    depends_on:
      - rabbitmq
      - postgres
    networks:
      - app-network
    restart: unless-stopped
    ports:
      - "${REPORT_PORT:-8081}:${REPORT_PORT:-8081}"
    environment:
      - DB_HOST=${DB_HOST}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_PORT=${DB_PORT}
      - DB_TIMEZONE=${DB_TIMEZONE}
      - RABBITMQ_HOST=${RABBITMQ_HOST}
      - RABBITMQ_USER=${RABBITMQ_USER}
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - RABBITMQ_PORT=${RABBITMQ_PORT}
      - REPORT_PORT=${REPORT_PORT}

volumes:
  rabbitmq_data:
  postgres_data:

networks:
  app-network:
    driver: bridge



#for development

services:
  rabbitmq:
    image: "rabbitmq:4.0.5-management-alpine"
    container_name: rabbitmq
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:15672/api/overview"]
      interval: 10s
      timeout: 5s
      retries: 5
    ports:
      - "5672:5672"
      - "15672:15672"
    networks:
      - app-network
    volumes:
      - "rabbitmq_data:/data"

  postgres:
    image: postgres:latest
    container_name: postgres
    environment:
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=${DB_NAME}
    ports:
      - "5432:5432"
    networks:
      - app-network
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:latest
    container_name: redis
    ports:
      - "6379:6379"
    networks:
      - app-network
    restart: unless-stopped

  api:
    build:
      context: api
      dockerfile: Dockerfile
      target: development
    depends_on:
      - rabbitmq
      - postgres
      - redis
    ports:
      - "${SERVER_PORT:-8080}:${SERVER_PORT:-8080}"
    networks:
      - app-network
    restart: unless-stopped
    environment:
      - DB_HOST=${DB_HOST}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_PORT=${DB_PORT}
      - DB_TIMEZONE=${DB_TIMEZONE}
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - RABBITMQ_HOST=${RABBITMQ_HOST}
      - RABBITMQ_USER=${RABBITMQ_USER}
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - RABBITMQ_PORT=${RABBITMQ_PORT}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - EMAIL_FROM=${EMAIL_FROM}
      - EMAIL_PASSWORD=${EMAIL_PASSWORD}
      - REPORT_PORT=${REPORT_PORT}
      - SERVER_PORT=${SERVER_PORT}
      - CASBIN_MODEL_PATH=/configs/model.conf
      - CASBIN_POLICY_PATH=/configs/policy.csv
    volumes:
      - ./api:/app
      - ./api/.air.toml:/app/.air.toml
      - ./api/pkg/auth/casbin/model.conf:/configs/model.conf
      - ./api/pkg/auth/casbin/policy.csv:/configs/policy.csv

  report:
    build:
      context: report
      dockerfile: Dockerfile
      target: development
    depends_on:
      - rabbitmq
      - postgres
    networks:
      - app-network
    restart: unless-stopped
    ports:
      - "${REPORT_PORT:-8081}:${REPORT_PORT:-8081}"
    environment:
      - DB_HOST=${DB_HOST}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_PORT=${DB_PORT}
      - DB_TIMEZONE=${DB_TIMEZONE}
      - RABBITMQ_HOST=${RABBITMQ_HOST}
      - RABBITMQ_USER=${RABBITMQ_USER}
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - RABBITMQ_PORT=${RABBITMQ_PORT}
      - REPORT_PORT=${REPORT_PORT}
    volumes:
      - ./report:/app

  mail:
    build:
      context: mail
      dockerfile: Dockerfile
      target: development
    depends_on:
      - rabbitmq
    networks:
      - app-network
    restart: unless-stopped
    environment:
      - RABBITMQ_HOST=${RABBITMQ_HOST}
      - RABBITMQ_USER=${RABBITMQ_USER}
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - RABBITMQ_PORT=${RABBITMQ_PORT}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - EMAIL_FROM=${EMAIL_FROM}
      - EMAIL_PASSWORD=${EMAIL_PASSWORD}
    volumes:
      - ./mail:/app

volumes:
  rabbitmq_data:
  postgres_data:

networks:
  app-network:
    driver: bridge

```

4. Create a .env file in the root directory of the project and add the following variables:

```bash
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_PORT=5432
DB_TIMEZONE=Asia/Jakarta
JWT_SECRET_KEY=secret
RABBITMQ_HOST=localhost
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_PORT=5672
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis
SMTP_HOST=localhost
SMTP_PORT=1025
EMAIL_FROM=noreply@example.com
EMAIL_PASSWORD=password
REPORT_PORT=8081
SERVER_PORT=8080
```

Directory Structure
```bash
├── docker compose.yml
├── go-super-farmer-api
├── go-super-farmer-report-service
├── go-super-farmer-mail-service
├── predict-model
└── .env
```

5. Run docker compose file
```bash
#for the first time build
docker compose up --build

#for runing all services
docker compose up
#or
docker compose start

#for stoping all services
docker compose down
#or
docker compose stop
```
