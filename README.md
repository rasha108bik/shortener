# go-musthave-shortener-tpl

Шаблон репозитория для практического трек "Веб-разработка на Go"

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команды:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

затем добавьте полученые изменения в свой репозиторий.

# Запуск автотестов

Для успешного запуска автотестов вам необходимо давать вашим веткам названия вида `iter<number>`, где `<number>` -
порядковый номер итерации.

Например в ветке с названием `iter4` запустятся автотесты для итераций с первой по четвертую.

При мерже ветки с итерацией в основную ветку (`main`) будут запускаться все автотесты.

# architecture
Using clean architecture https://github.com/evrone/go-clean-template
