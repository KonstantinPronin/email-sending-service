# email-sending-service

#### Description:
Сервис отправки почтовых уведомлений, состоящий из двух микросервисов: Acceptor и Sender.  
Acceptor принимает и обрабатывает HTTP-запросы в рамках API.  
Sender принимает от Acceptor-а уведомления, отправляет их через SMTP-сервер и сохраняет отправленные уведомления в базу.  
Acceptor и Sender взаимодействуют друг с другом через очередь. Когда Acceptor получает запрос на отправку уведомления, он пишет его в очередь.  
Sender вычитывает уведомление из очереди, устанавливает соединение с SMTP-сервером и отправляет полученное уведомление на указанный в нём адрес.  
После этого сохраняет отправленное уведомление в базу.

#### How to run:
Перед запуском необходимо настроить [конфигурационные файлы](./conf).

```shell script
    git clone https://github.com/KonstantinPronin/email-sending-service.git
    cd email-sending-service
    docker-compose build
    docker-compose up
```