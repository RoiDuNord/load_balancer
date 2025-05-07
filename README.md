## Технологии
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

## Предварительные требования

- Установленный Docker
- Базовые знания команд Docker

## Сборка Docker-образа и запуск приложения

1) Склонировать проект на ваш компьютер с Github с помощью команды:
```
git clone https://github.com/FilimonovAlexey/load_balancer.git
```

2) Билдинг образа
```
docker build -t lb -f docker/Dockerfile .
```

-t lb: задает имя для образа lb.
-f docker/Dockerfile: указывает путь к Dockerfile.
Запуск Docker-контейнера
После сборки образа вы можете запустить контейнер с помощью следующей команды:

3) Запуск контейнера
```
docker run -p 8080:8080 lb
```

4) Доступ к запущенному контейнеру
```
docker exec -it {container_id} /bin/sh // замените {container_id} на фактический ID контейнера
```

```
docker ps
```

4) Просмотр логов

Перейдите в директорию логирования:
```
cd logger
```

Используйте следующую команду, чтобы динамически отслеживать файл логов:
```
tail -f logger.log
```


![React](https://img.shields.io/badge/-React-61daf8?logo=react&logoColor=black)
![HTML5](https://img.shields.io/badge/-HTML5-e34f26?logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/-CSS3-1572b6?logo=css3&logoColor=white)
![JavaScript](https://img.shields.io/badge/-JavaScript-f7df1e?logo=javaScript&logoColor=black)
![Webpack](https://img.shields.io/badge/-Webpack-99d6f8?logo=webpack&logoColor=black)
