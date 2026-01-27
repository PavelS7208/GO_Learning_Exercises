## Реализация COW (Copy-On-Write) буфера

Буфер шарит данные между своими копиями пока значения не меняются. При создании клона увеличиваем счетчик. При изменении - создаем копию. Буфер следует закрывать после окончания использования


Реализованы методы

* func NewCOWBuffer(data []byte) COWBuffer - конструктор
* func (b *COWBuffer) Clone() COWBuffer  - Клонировать буфер
* func (b *COWBuffer) Close()      - закрыть буфер (уменьшить счетчик использования)
* func (b *COWBuffer) Update(index int, value byte) bool         - изменить указанный элемент
* func (b *COWBuffer) String() string  - преобразованеи в строку хранящегося набора byte (без копирования)


```go
   data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()
	
	copy1 := buffer.Clone()
	defer copy1.Close()
	copy2 := buffer.Clone()
	defer copy2.Close()

	fmt.Println(buffer.Update(0, 'g'))  // true

    fmt.Println(buffer.String())     // Измененная строка
	fmt.Println(copy1.String())      //  Исходная строка
	fmt.Println(copy2.String())      //  Исходная строка
```

