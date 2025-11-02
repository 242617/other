---
"[mdprc:skip_execute]": true
"[mdprc:skip_place]": true
"[mdprc:remove_properties]": true
---
## Ревью собственного кода

- Проводите review кода Go (веб-сервисы, микросервисы, CLI-утилиты, скрипты)
- Проверьте соблюдение Effective Go, Code Review Comments, информативные комментарии к экспортируемым функциям/типам
- Используйте линтеры: golangci-lint (включает staticcheck, govet, errcheck и др.)
- Проверяйте наличие unit/integration/table-driven тестов (testing, testify)
- Внимательно исследуйте обработку ошибок (error wrapping, sentinel errors), управление goroutines и контекстами (context.Context, sync.WaitGroup)
- Учитывайте best practices тестирования: table-driven tests, testify/suite, testify/mock, httptest
- Проверяйте корректное закрытие ресурсов через defer
- Обращайте внимание на race conditions (запускайте тесты с флагом -race)
- Пример замечания:
```go
// ⚠️ Ошибка безопасности: SQL-инъекция
query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", userID)
db.Query(query)
// 💬 Рекомендация: Использовать параметризованные запросы
db.Query("SELECT * FROM users WHERE id = $1", userID)
```

### Философия review:
- Конструктивная критика, а не просто указание на ошибки
- Объяснение "почему", а не только "что" нужно исправить
- Признание хороших решений наравне с указанием на проблемы
- Фокус на обучении и улучшении, а не на критике
- Evidence-based замечания (со ссылками на PRD, FR, стандарты, документацию)

### Процесс ревью

#### ШАГ 1: Подготовка к review

Перед началом review вы **ДОЛЖНЫ**:

```
CHECKLIST перед началом review:
□ Прочитан соответствующий раздел PRD (FR-XX)
□ Понятны все acceptance criteria из PRD, FR
□ Известны требования к качеству кода (coverage, linting)
□ Проверены связанные issue/tickets
□ Понятен контекст изменений (что и почему меняется)
□ Известны зависимости от других модулей (из PRD, FR)
```

Если **хотя бы один пункт НЕ выполнен** — запросите информацию:

```
❓ ЗАПРОС КОНТЕКСТА

Для проведения качественного review необходимо:

ОТСУТСТВУЕТ:
- Ссылка на PRD секцию FR-XX
- Описание проблемы, которую решает PR
- Контекст изменений

Пожалуйста, предоставьте:
1. Ссылку на соответствующий FR в PRD
2. Краткое описание проблемы, изменений, контекст
```

#### ШАГ 2: Функциональный review

Проверьте соответствие кода требованиям PRD:

**Критерии проверки:**

1. **Соответствие Acceptance Criteria**

Для каждого AC из PRD проверьте:
```
✅ AC1: [Критерий из PRD] — реализован в [file:line]
⚠️ AC2: [Критерий из PRD] — реализован частично, отсутствует [деталь]
❌ AC3: [Критерий из PRD] — не реализован
```

2. **Обработка Edge Cases**

Проверьте, что все граничные случаи из PRD обработаны:
```go
// ✅ ХОРОШО: Edge case из PRD обработан
func ProcessData(input string) error {
    // PRD Edge Case: пустой input
    if input == "" {
        return ErrEmptyInput
    }
    
    // PRD Edge Case: слишком длинный input
    if len(input) > MaxInputLength {
        return ErrInputTooLong
    }
    
    // основная логика...
    return process(input)
}

// ❌ ПЛОХО: Edge cases не обработаны
func ProcessData(input string) error {
    // сразу обработка без проверок
    return process(input)
}
```

3. **Обработка ошибок согласно PRD**

```
✅ Проверьте, что:
- Все error codes из PRD секции "Обработка ошибок" реализованы
- Ошибки возвращаются с правильными типами
- Error messages соответствуют PRD спецификации
- Есть proper wrapping для трассировки
```

#### ШАГ 3: Качество кода

##### Читаемость

###### Именование:
```go
// ✅ ХОРОШО: самодокументируемые имена
func CalculateUserDiscount(user User, order Order) (Discount, error) {
    discountRate := getDiscountRateByLoyaltyLevel(user.LoyaltyLevel)
    return applyDiscount(order.TotalAmount, discountRate), nil
}

// ❌ ПЛОХО: неясные имена
func calc(u User, o Order) (float64, error) {
    r := get(u.L)
    return apply(o.T, r), nil
}
```

**Замечание для плохого примера:**
```
💬 ПРЕДЛОЖЕНИЕ: Улучшить именование

Текущие имена функций и переменных слишком короткие и неясные:
- `calc` → `CalculateUserDiscount` (четко указывает, что вычисляется)
- `u`, `o` → `user`, `order` (в пределах функции допустимо)
- `r` → `discountRate` (ясно, что это за значение)

Почему: Код читается гораздо чаще, чем пишется. Явные имена делают код самодокументируемым и снижают когнитивную сложность.

Ссылка: Effective Go - Naming conventions
```

###### Комментарии
```go
// ✅ ХОРОШО: комментарии объясняют "почему", не "что"
// Используем exponential backoff т.к. внешний API имеет rate limiting
// и может временно отклонять запросы (из PRD NFR-05)
for attempt := 0; attempt < maxRetries; attempt++ {
    if err := callExternalAPI(); err == nil {
        return nil
    }
    time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
}

// ❌ ПЛОХО: комментарии дублируют код
// Цикл от 0 до maxRetries
for attempt := 0; attempt < maxRetries; attempt++ {
    // Вызов внешнего API
    callExternalAPI()
}
```

###### Функции — Single Responsibility
```go
// ✅ ХОРОШО: каждая функция делает одну вещь
func ProcessOrder(order Order) error {
    if err := validateOrder(order); err != nil {
        return errors.Wrap(err, "validate order")
    }
    
    if err := saveOrder(order); err != nil {
        return errors.Wrap(err, "save order")
    }
    
    if err := sendConfirmation(order); err != nil {
        return errors.Wrap(err, "send confirmation")
    }
    
    return nil
}

// ❌ ПЛОХО: функция делает слишком много
func ProcessOrder(order Order) error {
    // 200 строк валидации, сохранения, отправки email, 
    // логирования, метрик, обработки платежа...
}
```

##### Тестирование

Проверьте наличие тестов:
```
ОБЯЗАТЕЛЬНЫЕ тесты:
✅ Unit тесты для всех публичных функций
✅ Тесты для всех acceptance criteria из PRD
✅ Тесты для всех edge cases из PRD
✅ Тесты для error handling
✅ Coverage > 80% для критичных путей
```

**Пример замечания о тестах:**
```
⚠️ БЛОКЕР: Недостаточное покрытие тестами

Функция `CalculateDiscount` реализует критичную бизнес-логику из PRD FR-12, но отсутствуют тесты для:
1. Edge case: нулевая сумма заказа (упомянут в PRD)
2. Edge case: отрицательная скидка
3. Acceptance Criteria AC3: максимальная скидка не может превышать 50%

Требуется:
- Добавить тесты для всех edge cases из PRD
- Покрыть все acceptance criteria
- Достичь coverage > 80% для этой функции

Почему блокер: Без тестов невозможно гарантировать корректность реализации критичной бизнес-логики, что может привести к финансовым потерям.
```

##### Performance

###### N+1 проблемы

```go
// ❌ ПЛОХО: N+1 запросы к БД
func GetUsersWithOrders(userIDs []int) ([]UserWithOrders, error) {
    var result []UserWithOrders
    for _, id := range userIDs {
        user := getUserByID(id)  // N запросов
        orders := getOrdersByUserID(id)  // еще N запросов
        result = append(result, UserWithOrders{user, orders})
    }
    return result, nil
}

// ✅ ХОРОШО: batch загрузка
func GetUsersWithOrders(userIDs []int) ([]UserWithOrders, error) {
    users := getUsersByIDs(userIDs)  // 1 запрос
    orders := getOrdersByUserIDs(userIDs)  // 1 запрос
    return joinUsersAndOrders(users, orders), nil
}
```

###### Избыточные аллокации

```go
// ❌ ПЛОХО: аллокация в цикле
for i := 0; i < 1000; i++ {
    result := fmt.Sprintf("Item %d", i)  // 1000 аллокаций
    items = append(items, result)
}

// ✅ ХОРОШО: pre-allocate + strings.Builder
items := make([]string, 0, 1000)
var builder strings.Builder
for i := 0; i < 1000; i++ {
    builder.Reset()
    builder.WriteString("Item ")
    builder.WriteString(strconv.Itoa(i))
    items = append(items, builder.String())
}
```

**Замечание о performance:**
```
💡 ОПТИМИЗАЦИЯ: Можно улучшить производительность

Текущая реализация делает N+1 запросов к БД, что при большом количестве пользователей приведет к деградации производительности.

PRD NFR-02 требует: "поддержка до 1000 concurrent users"
Текущая реализация: при 1000 users = 2000 DB queries

Предложение:
- Использовать batch loading (getUsersByIDs, getOrdersByUserIDs)
- Сократит запросы с 2000 до 2
- Улучшит response time с ~2s до ~200ms (по бенчмаркам)

Референс: [ссылка на best practices или документацию]
```

##### Security

###### Input validation

```go
// ❌ ОПАСНО: нет валидации
func GetUserByID(id string) (*User, error) {
    query := "SELECT * FROM users WHERE id = " + id  // SQL injection!
    return db.Query(query)
}

// ✅ БЕЗОПАСНО: prepared statements
func GetUserByID(id string) (*User, error) {
    if !isValidUUID(id) {
        return nil, ErrInvalidID
    }
    query := "SELECT * FROM users WHERE id = $1"
    return db.Query(query, id)
}
```

###### Sensitive data handling

```go
// ❌ ОПАСНО: логирование паролей
log.Printf("User login: username=%s, password=%s", username, password)

// ✅ БЕЗОПАСНО: не логируем sensitive data
log.Printf("User login attempt: username=%s", username)
```

**Замечание о security:**
```
🚨 КРИТИЧНО: Уязвимость безопасности

Обнаружена потенциальная SQL injection уязвимость в функции `GetUserByID`.

Проблема:
- Прямая конкатенация user input в SQL запрос (line 45)
- Нет валидации входного параметра `id`

Риск:
- Attacker может выполнить произвольные SQL команды
- Возможна утечка данных или модификация БД

Решение:
1. Использовать prepared statements ($1, $2, etc.)
2. Добавить валидацию `id` (должен быть валидный UUID)
3. Добавить тест для SQL injection attempt

Ссылка: OWASP Top 10 - Injection
```
