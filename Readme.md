# REST API сервис объявлений

Этот репозиторий содержит REST API сервис для создания и просмотра объявлений. Сервис предоставляет несколько конечных точек для обработки различных функций, связанных с аутентификацией пользователей, регистрацией, размещением объявлений и отображением ленты объявлений.

## Конечные точки

1. **Аутентификация пользователя**
   - Конечная точка: `/login`
   - Метод: `POST`
   - Тело запроса: JSON с полями `login` и `password`
   - Полученный токен необходимо передавать в хедере `Authorization-access`

2. **Регистрация пользователя**
   - Конечная точка: `/register`
   - Метод: `POST`
   - Тело запроса: JSON с полями `login` и `password`

3. **Размещение объявления**
   - Конечная точка: `/advert`
   - Метод: `POST`
   - Тело запроса: JSON с полями `header`, `body`, `image_url` и `price`

4. **Отображение ленты объявлений**
   - Конечная точка: `/feed`
   - Метод: `GET`
   - Query parameters:
     - `sort`: Сортировка объявлений (`priceUp`, `priceDown`, `new`, `old`)
     - `priceMin`: Минимальная цена для фильтрации объявлений
     - `priceMax`: Максимальная цена для фильтрации объявлений
     - `page`: Номер страницы для пагинации

## Запуск Сервиса

### Использование Docker

1. Вы можете загрузить готовый Docker образ с Docker Hub
```bash
    docker pull yasminworks/admarket
```

2. Создание Docker образа
```bash
   docker build . -t admarket:latest
``` 

3. Создание Docker volume и запуск контейнера
```bash
   docker volume create ad-market
   docker run -d -it -p 8082:8082 -v ad-market:/app/storage admarket
```
   
### Компиляция Исходного Кода
0. Перейдите в корневую папку проекта

1. Установка зависимостей:
```bash
  go mod download
 ```

2. Подготовка базы данных:
```bash
   go run ./cmd/migrator --storage-path=./storage/storage.db --migrations-path=./migrations
 ```  
3. Компиляция и запуск:
 ```bash 
    go build -o ad-market ./cmd/ad-market/main.go
    ./ad-market
 ```  
### Использование Утилиты Task

Если у вас установлена утилита Task, можно запустить сервис командой
```bash
    task build
```
