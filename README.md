https://study.effective-mobile.ru/cources/go/stage_3/3-devops-docker-i-cicd/topic.html#практика

За основу взят код из прошлого задания, находится в internal\handlers\handlers.go

CI, прометей и графана в корневой папке, а так же в internal

## Запуск

```bash
docker compose up --build -d
```

## Проверка

```bash
curl "http://localhost:8080/healthz"
curl "http://localhost:8080/readyz"
curl "http://localhost:8080/cache/set?key=hello&value=world"
curl "http://localhost:8080/cache/get?key=hello"
curl "http://localhost:8080/users"
curl "http://localhost:8080/metrics"
```

```
Прометей: http://localhost:9090
Графана: http://localhost:3000
```