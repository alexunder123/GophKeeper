# GophKeeper

Основные тезисы работы программы:

Сервер хранит список зарегистрированных пользователей в SQL базе данных.
Пароль пользователя хранится в базе в зашифрованном виде.
При регистрации новому пользователю присваивается userID и генерируется симметричный ключ шифрования.
Симметричный ключ передается клиенту для зашифровки/расшифровки пользовательских данных.
Пользовательские данные клиент передает в зашифрованном виде, сервер сохраняет их в отдельном двоичном файле.

При загрузке клиентское приложение подключается к серверу запрашивает SessionId и обменивается с сервером открытыми ключами.

Дальнейший обмен данными производится в зашифрованном виде и проверкой подписей клиента и сервера.

При регистрации клиент отправляет логин и пароль, сервер шифрует и возвращает симметричный ключ и userID
При возникновении ошибки, процесс начинается заново

При авторизации клиент отправляет логин и пароль, сервер возвращает userID.

После авторизации клиент запрашивает данные по SessionId .
Сервер шифрует симметричный ключ и возвращает зашифрованные данные и ключ

Клиент может вручную проверять актуальный статус своих данных. Сервер возвращает время последнего сохранения данных и наличие текущей блокировки на изменения.

При начале добавления или редактирования данных клиент запрашивает блокировку от изменений данных на сервере.
Сервер возвращает подтверждение блокировки и до какого времени данные заблокированы или отказ и время до которого другой пользователь их заблокировал.

После добавления данных клиент шифрует их симметричным ключом и отправляет на сервер.
Сервер проверяет текущую блокировку и сравнивает время последнего сохранения данных сервера и клиента, и сохраняет данные или возвращает ошибку.

LogOut: сервер удаляет sessionID, клиент стирает данные в оперативной памяти и готов к новому входу в систему
