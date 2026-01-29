## Реализация OrderedMap

Идея упорядоченного словаря заключается в том, что он будет реализован на основе бинарного дерева поиска (BST). Дерево будет строиться только по ключам элементов, значения элементов при построении дерева неучитываются. Элементы с одинаковыми ключами в упорядоченном словаре хранить нельзя.


Реализовано с использованием джинериков. 

```go
OrderedMap[K cmp.Ordered, V any]
```


Ограничения и особенности:
* Реализация простая (обычное дерево двоичное) для примера, может вырождаться в список
* Только базовые функции для проверкм
* Покрыто тестами



Реализованы методы

* func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V]  - конструктор
* func (m *OrderedMap[K, V]) Get(key K) (V, bool)  - Получение элемента по ключу
* func (m *OrderedMap[K, V]) Insert(key K, value V) bool - Вставить элемент (с проверкой что вставка произошла). ПРи существовании ключа - перезапись значения
* func (m *OrderedMap[K, V]) Delete(key K) bool - Удаление по ключу
* func (m *OrderedMap[K, V]) Contains(key K) bool - Проверка на наличие по ключу
* func (m *OrderedMap[K, V]) Size() int / func (m *OrderedMap[K, V]) Empty() bool - Размеры
* func (m *OrderedMap[K, V]) ForEachInOrder(action func(K, V))   - Итератор forEach в порядке возрастания ключей
* func (m *OrderedMap[K, V]) String() string  - преобразованеи в строку хранящихся значений (для тестирования)




```go
	data := NewOrderedMap[int, int]()
	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	fmt.Println(data.String())
  // Обход от наименьшего ключа к наибольшему
	var index int
	size := data.Size()
	data.ForEachInOrder(func(k, v int) {
		switch index {
		case 0:
			fmt.Println("Наименьший элемент", v)
		case size - 1:
			fmt.Println("Наибольший элемент", v)
		}
		index++
	})


```
