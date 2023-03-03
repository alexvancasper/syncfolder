# Synchronizer
Synchronizer синхронизирует 2 папки между собой включая права на файлы.
Имеется конфигурационный файл в yaml формате с комментариями для каждого поля.
Может работать в фоне в *nix среде с использования амперсанда. Пример показан ниже.

## Compilation
```sh
git clone <repository path>
cd <repository>
```

To run tests
```
make test
```

To build binary file
```
make build
```
In the folder will appear new file <repository>/cmd/app/syncher which is ready to use.

## Usage
Run in the background
```sh
$ cmd/app/syncher config/config.yml &
```
Stop it
```sh
$ pkill syncher

```

