package ozon

import (
	"fmt"
	"math"
	"sort"
  )  

//На Авито размещено множество товаров, каждый из которых представлен числом. 
//У каждого покупателя есть потребность в товаре, также выраженная числом. 
//Если точного товара нет, покупатель выбирает ближайший по значению товар, 
//что вызывает неудовлетворённость, равную разнице между его потребностью и купленным товаром. 
//Количество каждого товара не ограничено, и один товар могут купить несколько покупателей. 
//Рассчитайте суммарную неудовлетворённость всех покупателей.

//Нужно написать функцию, которая примет на вход два массива: 
//массив товаров и массив потребностей покупателей, вычислит сумму неудовлетворённостей всех покупателей и вернет результат в виде числа.

//Пример
//# ввод
//goods = [8, 3, 5]
//buyerNeeds = [5, 6]
//# вывод
//res = 1 # первый покупатель покупает товар 5 и его неудовлетворённость = 0, второй также покупает товар 5 и его неудовлетворённость = 6-5 = 1


// функция для поиска ближайшего товара
func findClosest(goods []int, need int) int {
  // используем двоичный поиск для нахождения места
  idx := sort.Search(len(goods), func(i int) bool {
    return goods[i] >= need
  })

  // обрабатываем граничные случаи
  if idx == 0 {
    return goods[0]
  }
  if idx == len(goods) {
    return goods[len(goods)-1]
  }

  // выбираем ближайший товар из двух возможных
  if math.Abs(float64(goods[idx]-need)) < math.Abs(float64(goods[idx-1]-need)) {
    return goods[idx]
  }
  return goods[idx-1]
}

func calculateTotalDissatisfaction(goods []int, buyerNeeds []int) int {
  // сортируем товары для эффективного поиска
  sort.Ints(goods)
  totalDissatisfaction := 0

  // проходим по всем покупателям
  for _, need := range buyerNeeds {
    closest := findClosest(goods, need)
    dissatisfaction := int(math.Abs(float64(closest - need)))
    totalDissatisfaction += dissatisfaction
  }

  return totalDissatisfaction
}

func T1() {
  goods := []int{8, 3, 5}
  buyerNeeds := []int{5, 6}
  result := calculateTotalDissatisfaction(goods, buyerNeeds)
  fmt.Println(result) // ожидается 1
}