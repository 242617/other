---
"[mdprc:skip_execute]": true
"[mdprc:skip_place]": true
"[mdprc:remove_properties]": true
---
**ШАГ 3: GREEN — Реализация минимального кода**

Теперь пишите **минимальный** код для прохождения тестов:

**Правила идиоматичного Go:**

1. ЯВНАЯ обработка ошибок с контекстом:
```go
func ProcessData(data InputType) (*OutputType, error) {
    if err := validateInput(data); err != nil {
        return nil, errors.Wrap(err, "validate input")
    }
    
    result, err := performOperation(data)
    if err != nil {
        return nil, errors.Wrap(err, "perform operation")
    }
    
    return result, nil
}
```

2. Используйте интерфейсы для абстракций (где уместно):
```go
type DataProcessor interface {
    Process(InputType) (OutputType, error)
}
```

3. Короткие, фокусированные функции
Если функция > 50 строк — разбейте на подфункции

4. Идиоматичные имена переменных
ctx - context, err - error, cfg - config, db - database
Избегайте длинных имен: используйте короткие в ограниченном scope

5. Используйте defer для cleanup
```go
func ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return errors.Wrap(err, "os open")
    }
    defer file.Close() // гарантирует закрытие
    
    // работа с файлом...
    return nil
}
```

6. **Структура пакетов и именование:**
   - Используйте описательные имена пакетов (строчные буквы, без подчеркиваний)
   - Используйте camelCase для переменных и функций, PascalCase для экспортируемых типов
   - Группируйте импорты по стандартной библиотеке, внешним зависимостям и внутренним пакетам

7. **Организация кода:**
   - Размещайте определения типов перед функциями в файлах
   - Группируйте связанные функции вместе
   - Используйте последовательный порядок полей структур: публичные поля, внутренние поля, логгер

8. **Логирование:**
   - Используйте структурированное логирование с парой ключ-значение:
   ```go
   log.Info(ctx, "message", "key", value)
   ```
   - Включайте соответствующий контекст в сообщения лога
   - Записывайте ошибки с информацией о модуле, методе и параметрах

9. **Архитектура сервисов:**
   - Используйте внедрение зависимостей для сервисов
   - Используйте паттерн модификаторов/опций для создания и настройки структур
   - Используйте интерфейсы инструментов для расширяемости

10. **Использование контекста:**
    - Всегда передавайте контекст в функции, которые могут быть отменены или истечь по времени

11. **Конфигурация:**
    - Предпочтительный формат конфигурации – YAML
    - Используйте константы для значений, которые не должны изменяться

12. **Валидация:**
    - Реализуйте методы валидации для структур данных
    - Используйте специфичные типы ошибок для разных сценариев ошибок

**Обработка ошибок согласно PRD:**
- Каждый тип ошибки из PRD должен быть именованным: `var ErrInvalidInput = errors.New("invalid input")`
- Используйте `errors.Wrap` с именем компонента и его метода
- Возвращайте zero value + error для types
- Возвращайте nil + error для pointer types

**Concurrency:**
```go
// Используйте sync.WaitGroup для coordination
var wg sync.WaitGroup
results := make(chan Result, len(items))

for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        result := processItem(item)
        results <- result
    }(item) // передаём копию для избежания race condition
}

wg.Wait()
close(results)
```
