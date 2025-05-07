# Load Balancer

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

## Описание
Этот проект представляет собой простой балансировщик нагрузки, написанный на Go и упакованный в Docker-контейнер. Он позволяет распределять входящие запросы между несколькими серверами.

## Предварительные требования
- **Docker**: Убедитесь, что Docker установлен и работает на вашем компьютере. Команда выведет установленную версию Docker. Если Docker не установлен, вы получите сообщение об ошибке.
```
docker --version
```
- проверка работы Docker
```
docker info
```
- если Docker работает, вы увидите информацию о конфигурации и состоянии вашего Docker-демона. Если он не запущен, вы получите сообщение об ошибке
```
docker info
```

## Алгоритм сборки Docker-образа и запуска приложения

1) Клонировать проект на ваш компьютер с Github с помощью команды
```
git clone https://github.com/FilimonovAlexey/load_balancer.git
```

2) Билдинг образа
```
docker build -t lb -f docker/Dockerfile .
```
- -t lb: задает имя для образа lb
- -f docker/Dockerfile: указывает путь к Dockerfile

3) Запуск контейнера
```
docker run -p 8080:8080 lb
```

4) Доступ к запущенному контейнеру
- получить фактический ID контейнера
```
docker ps
```
- замените {container_id} на фактический ID контейнера
```
docker exec -it {container_id} /bin/sh
```

4) Просмотр логов
- перейдите в директорию логирования
```
cd logger
```
- команда для динамического отслежививания логов
```
tail -f logger.log
```
