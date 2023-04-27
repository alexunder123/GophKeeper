# GophKeeper

план:

сервер хранит базу клиентов в чем?
данные клиента хранятся в структуре JSON в отдельных файлах в зашифрованном виде

сервер хранит и передает данные клиентов в зашифрованном виде AES + chiper для деления по блокам

При регистрации клиент запрашивает открытый ключ, сервер возвращает ключ и SessionId
клиент шифрует логин и пароль и отправляет запрос. В ответ приходит или userID и ОК или ошибка "логин занят"
При возникновении ошибки, процесс начинается заново

При авторизации клиент запрашивает открытый ключ, сервер возвращает ключ и SessionId.
Клиент шифрует логин и пароль и отправляет запрос. В ответ приходит userID

После авторизации клиент запрашивает данные по userID и отправляет открытый ключ.
Сервер шифрует симметричный ключ и возвращает зашифрованные данные и ключ

После добавления данных клиент шифрует их открытым ключом и отправляет на сервер.
Сервер сравнивает время последнего сохранения данных сервера и клиента, и сохраняет данные или возвращает ошибку "данные клиента не актуальны"

После обновления данных на сервере, он инициализирует процесс обновления данных у клиентов???
1. Либо клиент периодически делает запросы к серверу для проверки отметки времени изменений
2. Либо при входе в редактирование данных делать предварительно запрос, и если они изменились предложить обновить их
Либо комбинация вариантов 1 и 2

LogOut: сервер удаляет sessionID, клиент стирает данные в оперативной памяти и готов к новому входу в систему

Сервер хранит оперативные данные в мапе где ключ это SessionId, а данные хранятся в структуре {userID, OpenKey, PrivateKey, Expires}
<!-- Временные идентификаторы сервер хранит в мапе где ключ это userID, а значение SessionId -->

Клиент хранит оперативные данные в структуре, {OpenKey, PrivateKey, Expires, SessionId, userID, ServerOpenKey, AES-ключ, данные клиента{}}

Нужно ли подписывать зашифрованные сообщения?

меню пользователя в бесконечном цикле через switch вызываем функции подменю