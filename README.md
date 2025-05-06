
## Запуск

1. Создать config.yaml файл в корневом каталоге и добавьте следующие значения:

    ```yaml
   server:
      host: "rate-balancer"
      port: 8080
   
   adminServer:
      host: "rate-balancer"
      port: 8888

   database:
      postgres_db: "rateLimiter"
      postgres_host: "postgres"
      postgres_ports: "5432"
      postgres_user: "pguser"
      postgres_password: "pgpwd"

   loadBalancer:
      strategy: "round_robin"
      backends:
         - url: ""
         - url: ""
         - url: ""
      unhealthy_threshold: 
      healthy_threshold: 
      timeout: 

   healthChecker:
      ping_interval: 

   rateLimiter:
      default_capacity: 100
      default_rate_per_sec: 10
    ```

2. Запустить проект

    ```
   make up
   ```
   
3. Прогнать миграции
    ```
   make migrate
   ```
## Запуск тестов

1. Создать тестовую базу

    ```
   make test-db
   ```

2. Прогнать миграции
    ```
   make test-migrate
   ```
3. Запустить тесты
   ```
   make test
   ```