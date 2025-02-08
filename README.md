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

x-minio-common: &minio-common
  image: quay.io/minio/minio:RELEASE.2025-01-20T14-49-07Z
  command: server --console-address ":9001" http://minio{1...4}/data{1...2}
  networks:
    - app-network
  expose:
    - "9000"
    - "9001"
  environment:
    MINIO_ROOT_USER: ryvasa
    MINIO_ROOT_PASSWORD: Cobacoba
  healthcheck:
    test: ["CMD", "mc", "ready", "local"]
    interval: 5s
    timeout: 5s
    retries: 5

services:
  minio1:
    <<: *minio-common
    hostname: minio1
    volumes:
      - data1-1:/data1
      - data1-2:/data2

  minio2:
    <<: *minio-common
    hostname: minio2
    volumes:
      - data2-1:/data1
      - data2-2:/data2

  minio3:
    <<: *minio-common
    hostname: minio3
    volumes:
      - data3-1:/data1
      - data3-2:/data2

  minio4:
    <<: *minio-common
    hostname: minio4
    volumes:
      - data4-1:/data1
      - data4-2:/data2

  nginx:
    image: nginx:1.19.2-alpine
    hostname: nginx
    networks:
      - app-network
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "9000:9000"
      - "9001:9001"
    depends_on:
      - minio1
      - minio2
      - minio3
      - minio4
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

  report:
    build:
      context: report
      dockerfile: Dockerfile
      target: runner
    depends_on:
      - postgres
      - minio1
      - minio2
      - minio3
      - minio4
    networks:
      - app-network
    restart: unless-stopped
    ports:
      - "${REPORT_PORT:-50051}:${REPORT_PORT:-50051}"
    environment:
      - DB_HOST=${DB_HOST}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_PORT=${DB_PORT}
      - DB_TIMEZONE=${DB_TIMEZONE}
      - MINIO_ENDPOINT=${MINIO_ENDPOINT}
      - MINIO_ID=${MINIO_ID}
      - MINIO_SECRET=${MINIO_SECRET}

  api:
    build:
      context: api
      dockerfile: Dockerfile
      target: runner
    depends_on:
      - rabbitmq
      - postgres
      - redis
      - report
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
      - REPORT_SERVICE_HOST=${REPORT_SERVICE_HOST}
      - REPORT_SERVICE_PORT=${REPORT_SERVICE_PORT}
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
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - EMAIL_FROM=${EMAIL_FROM}
      - EMAIL_PASSWORD=${EMAIL_PASSWORD}

volumes:
  rabbitmq_data:
  postgres_data:
  data1-1:
  data1-2:
  data2-1:
  data2-2:
  data3-1:
  data3-2:
  data4-1:
  data4-2:

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
DB_PORT=5432
DB_NAME=go_super_farmer
DB_TEST=go_super_farmer_test
DB_USER=postgres
DB_PASSWORD=123
DB_TIMEZONE=Asia/Jakarta

RABBITMQ_PORT=5672
RABBITMQ_PASSWORD=guest
RABBITMQ_USER=guest

SERVER_PORT=8080
REPORT_PORT=50051

JWT_SECRET_KEY=a1dQ3@k3lew#Ev5vcCx%z6^z0Nnm8T)u3rMed3s*cs2%w2e1fH7t#uGii8u0o@oSo4o!oewc7+q9eT30u978Y

REDIS_PORT=6379
REDIS_PASSWORD=

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
EMAIL_FROM=oktaviandua4@gmail.com
EMAIL_PASSWORD=lfxdohvgydnklurr

CASBIN_MODEL_PATH=./pkg/auth/casbin/model.conf
CASBIN_POLICY_PATH=./pkg/auth/casbin/policy.csv

REPORT_SERVICE_HOST=report
REPORT_SERVICE_PORT=50051

MINIO_ENDPOINT=minio1:9000
MINIO_ID=ryvasa
MINIO_SECRET=Cobacoba

# Container
DB_HOST=postgres
REDIS_HOST=redis
RABBITMQ_HOST=rabbitmq

# # localhost
# DB_HOST=localhost
# REDIS_HOST=localhost
# RABBITMQ_HOST=localhost
# MINIO_ENDPOINT=localhost:9000

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
