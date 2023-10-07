FROM postgres:10.0-alpine

ENV POSTGRES_USER admin
ENV POSTGRES_PASSWORD admin
ENV POSTGRES_DB comment_service

# Копирование файла schema.sql внутрь контейнера
COPY db/schema.sql /docker-entrypoint-initdb.d/

# Установка правильных разрешений на файл
RUN chmod 755 /docker-entrypoint-initdb.d/schema.sql

# Открытие порта для доступа к PostgreSQL
EXPOSE 5432

# Запуск PostgreSQL при запуске контейнера
CMD ["postgres"]
