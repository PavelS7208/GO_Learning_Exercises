## Реализация кольцевой очереди (Circular Queue)


Использованы дженерики с ограничением на нумерик типы

```go
type CircularQueue[T Number] struct {
    values []T
    cap    int
    front  int
    size   int
}
```

Реализованы методы

* NewCircularQueue - конструктор
* Push(value) bool  - Вставить значение в буфер
* Pop() bool        - удалить первый поставленный элемент (FIFO)
* Front() T         - считать последний записанный элемент
* Back() T          - считать последний записанный элемент
* String() string   - преобразованеи в строку (для отладки и вывода)




