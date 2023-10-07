FROM postgres:10.0-alpine

# Открытие порта для доступа к PostgreSQL
EXPOSE 5432

# Запуск PostgreSQL при запуске контейнера
CMD ["postgres"]