# Лаб 1: порождающие паттерны

Сравнительный анализ приложения **с** и **без** использования паттерны проектирования Singleton.

---

## Запустить

```bash
docker compose up --build
```

откройте:

- **http://localhost:3000/no-pattern** — No Pattern
- **http://localhost:3000/pattern** — Singleton Pattern

---

## Разница

### Без паттерн (`no-pattern`)

```go
// each handler opens & closes its own conn:
func getTasks() {
    db := sql.Open("sqlite3", "tasks.db")  // new conn
    defer db.Close()
    // ... query ...
}

func createTask() {
    db := sql.Open("sqlite3", "tasks.db")  // again new conn
    defer db.Close()
    // ...
}
```

**Проблемы:**

- Расточительное использование ресурсов = накладные расходы на соединение при каждом запросе.
- Нет пула соединений
- Неотслеживаемое количество соединений
- Risk of "too many open files" under load

---

### Singleton паттерн (`pattern`)

```go
// 1 instance, shared via sync.Once:
var (
    instance *Database
    once     sync.Once
)

func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{conn: openDB()}  // only 1
    })
    return instance  // always the same object
}

func getTasks() {
    db := database.GetInstance().GetConn()  // use existing connection
}
```

**Преимущества:**

- Единый пул подключений - эффективный и предсказуемый.
- Потокобезопасное обеспечение с помощью `sync.Once`
- Централизованное управление жизненным циклом
- Отслеживаемая статистика (открытые соединения, количество звонков)
