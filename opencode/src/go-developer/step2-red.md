---
"[mdprc:skip_execute]": true
"[mdprc:skip_place]": true
"[mdprc:remove_properties]": true
---
**ШАГ 2: RED — Написание failing теста**

**ОБЯЗАТЕЛЬНАЯ последовательность:**

1. Создайте тестовый файл `<module>_test.go`
2. Напишите тест, описывающий **ЧТО** должна делать функция (не КАК)
3. Включите все acceptance criteria из PRD как test cases
4. Включите все edge cases из секции "Граничные случаи" PRD
5. Запустите тест — он **ДОЛЖЕН** упасть с понятной ошибкой

**Формат теста:**

```go
package module_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestFunctionName_Scenario описывает конкретный сценарий
// из PRD FR-XX Acceptance Criteria #1
func TestFunctionName_BasicScenario(t *testing.T) {
    // Arrange - подготовка входных данных из PRD
    input := module.InputType{
        Field: "value from PRD example",
    }
    
    // Act - вызов функции
    result, err := module.FunctionName(input)
    
    // Assert - проверка согласно acceptance criteria
    require.NoError(t, err, "PRD FR-XX не должна возвращать ошибку в базовом сценарии")
    assert.Equal(t, expectedOutput, result, "Результат должен соответствовать PRD спецификации")
}

// TestFunctionName_EdgeCase тестирует граничный случай
// из PRD секция "Граничные случаи"
func TestFunctionName_EmptyInput(t *testing.T) {
    // Тест для edge case...
}

// TestFunctionName_ErrorHandling тестирует обработку ошибок
// согласно PRD секция "Обработка ошибок"
func TestFunctionName_InvalidInput(t *testing.T) {
    // Arrange
    invalidInput := module.InputType{} // некорректные данные
    
    // Act
    _, err := module.FunctionName(invalidInput)
    
    // Assert - проверяем ожидаемый тип ошибки из PRD
    require.Error(t, err)
    assert.Contains(t, err.Error(), "expected error message from PRD")
}
```

**Критически важно:**
- Каждый acceptance criterion из PRD = отдельный test case
- Используйте table-driven tests для множественных сценариев
- Имена тестов должны явно указывать на секцию PRD
