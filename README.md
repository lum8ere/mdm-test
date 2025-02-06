Общее
----------------
1) Не стал делать автомиграции, но подготовил запросы
2) Создать любово пользователя, пароль сразу вставить зашифрованный (https://bcrypt.online/ - дефолтные настройки). И проставить роль админ или юзер (admin/user).
3) ну а тут уже можно баловаться через ui, либо через консоль
4) добавил make команды для запуска бека и клиента (андройд телефона)

ENV
----------------

Нужно создать **.env** файл в папке backend с такими полями:

```text
    DATABASE_URL="postgresql://postgres:root@localhost:5432/postgres?statement_timeout=120000"
    BC_PORT=:4000
    JWT_SECRET="your-secret"
```

Примеры запросов
----------------

Используйте инструменты вроде **Postman** или **cURL** для тестирования API.

-   **Получение токена:**

    ```bash
      curl -X POST http://localhost:4000/login \
      -H "Content-Type: application/json" \
      -d '{"username": "USER", "password": "PASSWORD"}'
    ```

-   **Изменение режима камеры:**

    ```bash
    curl -X POST http://localhost:4000/devices/android-test/camera \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
      -d '{"enabled": false}'
    ```