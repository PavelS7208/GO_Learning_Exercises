Реализовация кольцевой очереди (Circular Queue)


Использованы дженерики с ограничением на нумерик типы


type CircularQueue[T Number] struct {\
    values []T\
    cap    int\
    front  int\
    size   int\
}


