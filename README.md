# Locked

Это учебный проект для изучения написания полноценных приложений на языке Golang, изучения различных форматов файлов и изучения криптографии. 

Интерфейсы:

1. CLI
2. Web (local)
3. Возможно, даже desktop UI с помощью Fyne.

## Что из себя представляет проект?

Это хранилище секретов. Хранение будет в файлах собственного формата. Данные будут зашифрованы и сжаты, что сделает невозможным прочтение содержимого файлов без ключа.

3 основных пакета:
1. `cmd/locked`
2. `internal/core`
3. `internal/app`

`cmd/locked` будет располагать в себе `main`, где мы будем запускать просто приложуху.

В `internal/core` будет основная логика - шифрование, дешифрование и тд. ВСе это будет в `core`.

В `internal/app` будет код, связанный с обработкой команд. Да и в целом с вызовом методов из `core`. 

## Примерная схема файла:

1. Header - часть файла с метаданными, такими как версия, соль от пароля, контрольная сумма, количеством секретов и тд
2. Основная часть - сами секреты.

Чтобы по секретам можно было быстро переходить, у каждого секрета будет Offset, который будет указывать на следующий секрет.

## Проблемы

- Основная проблема: удаление секретов из начала/середины файла.
- Еще одна проблема: конфиденциальность. Да, это будет проблемой, так как мы не можем руками очистить определенный участок памяти, а GC не гарантирует, что очистка кучи произойдет сразу после взаимодействия с ней. Поэтому, скорее всего, переменные будут локальными и не указателями. К тому же нельзя будет держать файл открытым длительное время. Надо сделать так, чтобы каждый запрос мы открывали файл, выполняли какие-то действия и после сразу закрывали.

Предполагается использование **memguard**.

## Что можно хранить?

Все. Буквально все. Любой файл до определенного размера или просто текст - все можно будет хранить в таком хранилище.

## Структура секрета

Каждый секрет надо представить в виде какой-то структуры. Как это сделать?

Ну, 100% секрет должен иметь payload/полезную нагрузку/сами данные. Еще нужен `offset`, так как планируется сделать обход секретов именно через оффсеты. Это позволит как быстро искать все секреты через горутины/каналы, так и в целом распознавать программе, с какого момента начинается именно секрет.

Так, окей, а что еще? Еще нужно какое-то название, верно? Верно. Причем название должно быть ФИКСИРОВАННОЙ ДЛИНЫ, как и все поля, кроме данных =)

Еще в идеале выделить часть на описание формата файла, правильно? Да, надо бы, чтобы понимать, как в целом экспортировать секрет.

Можно добавить дату добавления секрета для красивого вывода и доп инфы. 

Еще надо добавить какое-то описание, чтобы можно было понять, че за секрет вообще.

Еще можно какой-нибудь ID типа int добавить.

Окей, вот что получаем:
1. **Offset** *(8 bytes)*
2. **ID** *(2-4 bytes)*
3. **Name** *(32-64 bytes)*
4. **Type** *(1 byte - будем хранить как число от 0 до 255)*
5. **Date** *(8 byte Unix)*
6. **Description** *(64-128 bytes)*
7. **Payload size** *(8 bytes для огромных секретов/файлов)*
8. **Payload** *(n bytes)*

## Жизненный цикл программы

Вообще изначально планировалось сделать так, чтобы приложение как бы и не было запущено. То есть мы вот в CLI вводим команду, приложение запускается, выполняет команду и завершается. Это помогло бы в теории предоставить большую защиту, однако есть минус - **добавлять секреты будет неудобно**. Нет, конечно, так программа сможет работать, однако хочется добавить нормальный CLI интерфейс.

К тому же надо будет сделать хороший WEB. Не знаю пока, что буду использовать, но это должно быть и красиво, и удобно. + Как-то добавить интеграцию с браузером, чтобы пароли можно было по наажтию на кнопку сохранять. 

Предполагаю, что при запуске в консоли и в WEB надо будет сначала ввести пароль (авторизоваться по сути). В этот момент будет запускаться приложение. Оно будет работать и ожидать ввода команд.

## Формат хранения

Крч, я не буду сжимать/шифровать ВСЮ инфу о секрете. Я буду это делать только с данными непосредственно.

**Какой плюс?** Сможем читать файл НЕ затрагивая чувствительные данные. Мы сможем выдать юзеру инфу о секретах всех, но там не будет самих секретов. Таким образом файл в целом можно открытым держать, пока юзер сам не закроет его, так как в памяти чувствительные данные будут лежать в сжатом+зашифрованном виде. Думаю, что это круто и правильно.

К тому же, это поможет читать содержимое максимально быстро. Та и нам не надо, чтобы юзер видел сразу все секреты. Я думаю, что и юзеру это не надо. Ему нужно зайти и посмотреть что-то конкретное, но чтобы все остальное было скрыто. 

Еще тогда можно будет дешифровать определенные запрашиваемые данные, тогда остальные секреты будут в безопасности.

## Примерная структура работы приложения

**2024.10.23**

Пускй пока что README.md будет чем-то по типу заметок на время.

Итак, я долго думал, как организовать общение всех слоев. 

Такая структура сейчас на примете:

UI (CLI, Web, Fyne mb) -> App -> WorkerPool -> Service -> Database.

То есть есть, условно, CLI UI. Пользователь вводит какую-то команду. Это прослушивается на уровне App, далее создается таска (worker.Task) и отправляется в канал. Канал читается в горутине на стороне WorkerPool. Какой-либо из воркеров обрабатывает эту таску, вызывая сервис и тд. Результат на уровне сервиса отправляется в канал resCh, который будет прослушиваться на уровне App (думаю, нет смысла слушать его на уровне WorkerPool). 