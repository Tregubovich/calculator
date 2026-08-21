# Calculator

Автор: Трегубович Андрей

## Установка и запуск
### Установка
```bash
git clone https://github.com/Tregubovich/calculator.git
cd calculator
go build -o out/calculator calculator/cmd
chmod u+x out/calculator
```

### Запуск
```bash
./out/calculator
```

## Поддерживаемые интерфейсы
- CLI
- HTTP

## Поддерживаемые операции
- `+-*/^`
- Константы: `pi`, `e`, `phi`
- Функции: `log2`, `log10`, `ln`, `sqrt`, `cbrt`, `abs`
- Тригонометрические функции: `sin`, `cos`, `tan`, `asin`, `acos`, `atan`,
- Факториал: `!`