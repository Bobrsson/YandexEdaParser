FROM golang:1.18.0 as base
FROM base as dev
# определяем корневую директорию. Текущую...
WORKDIR /home/vsvolkov/GolandProjects/YandexEdaParser
# копируем зависимости
COPY go.mod .
COPY go.sum .
#RUN go mod download
#копируем все файлики проекта в контейнер
COPY . $WORKDIR
# Собираем приложение, называем сборку
#RUN go build migrate up  -m ./migrations/ -o /migrator
RUN go build -o /YandexEdaParser
EXPOSE 8000
# вызываем только созданную сборку
CMD [ "/YandexEdaParser" ]

