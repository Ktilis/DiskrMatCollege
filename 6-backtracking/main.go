package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Ktilis/DiskrMat/common"
)

func main() {
	fmt.Println("DiskrMat: Поиск в глубину (Backtracking/DFS)")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nКакую задачу выполнить?")
	fmt.Println("2 - Задача 2: Подмножества (Все комбинации)")
	fmt.Println("3 - Задача 3: Путь в лабиринте (От начала до конца)")
	fmt.Println("4 - Задача 4: Генерация скобок")
	fmt.Println("5 - Задача 5: Комбинации сумм")
	fmt.Print("Ваш выбор: ")

	taskChoice := common.ReadLine(reader)

	fmt.Println("-------------------------")
	switch taskChoice {
	case "2":
		solveTask2()
	case "3":
		solveTask3()
	case "4":
		solveTask4()
	case "5":
		solveTask5()
	default:
		fmt.Println("Неизвестная задача.")
	}
}

func solveTask2() {
	fmt.Println("Решение Задачи 2: Подмножества")
	nums := []int{1, 2, 3}
	fmt.Printf("Вход: %v\n", nums)

	var result [][]int
	var current []int

	var backtrack func(start int)
	backtrack = func(start int) {
		// Делаем копию изначального массива, так как будем его изменять при возвратах
		temp := make([]int, len(current))
		copy(temp, current)
		result = append(result, temp)

		// Перебираем все возможные следующие элементы
		for i := start; i < len(nums); i++ {
			// Шаг вперед (добавляем элемент)
			current = append(current, nums[i])
			// Уходим в глубину
			backtrack(i + 1)
			// Откат/Backtracking (удаляем последний элемент)
			current = current[:len(current)-1]
		}
	}

	// Запускаем алгоритм с нулевого индекса
	backtrack(0)

	fmt.Printf("Выход: %v\n", result)
}

func solveTask3() {
	fmt.Println("Решение Задачи 3: Путь в лабиринте")

	maze := [][]int{
		{0, 0, 0},
		{1, 1, 0},
		{0, 0, 0},
	}

	rows := len(maze)
	cols := len(maze[0])

	fmt.Println("Вход (Лабиринт):")
	for _, row := range maze {
		fmt.Println(row)
	}

	var dfs func(r, c int) bool
	dfs = func(r, c int) bool {
		// 1. Проверка на выход за границы
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return false
		}

		// 2. Проверка на препятствие (1) или уже посещенную клетку (-1)
		if maze[r][c] == 1 || maze[r][c] == -1 {
			return false
		}

		// 3. Базовый случай успешного завершения: достигли правого нижнего угла
		if r == rows-1 && c == cols-1 {
			return true
		}

		// Помечаем клетку как посещенную
		maze[r][c] = -1

		// Идем в 4 направлениях. Если хоть где-то найдем путь - вернем true
		foundPath := dfs(r+1, c) || // Вниз
			dfs(r-1, c) || // Вверх
			dfs(r, c+1) || // Вправо
			dfs(r, c-1) // Влево

		return foundPath
	}

	hasPath := dfs(0, 0)

	if hasPath {
		fmt.Println("Существует путь: Да")
	} else {
		fmt.Println("Существует путь: Нет")
	}
}

func solveTask4() {
	fmt.Println("Решение Задачи 4: Генерация скобок")
	n := 3
	fmt.Printf("Вход: N=%d\n", n)

	var result []string

	var backtrack func(current string, openCount, closeCount int)
	backtrack = func(current string, openCount, closeCount int) {
		// Если длина строки достигла 2*N, мы сформировали правильную комбинацию
		if len(current) == n*2 {
			result = append(result, current)
			return
		}

		// Правило 1: Мы можем добавлять открывающую скобку, пока их меньше N
		if openCount < n {
			backtrack(current+"(", openCount+1, closeCount)
		}

		// Правило 2: Мы можем добавлять закрывающую скобку, ТОЛЬКО если
		// уже поставленных открывающих скобок больше, чем закрывающих
		if closeCount < openCount {
			backtrack(current+")", openCount, closeCount+1)
		}
	}

	backtrack("", 0, 0)

	fmt.Printf("Выход: %q\n", result)
}

func solveTask5() {
	fmt.Println("Решение Задачи 5: Комбинации сумм")
	nums := []int{2, 3, 6, 7}
	target := 7
	fmt.Printf("Вход: массив %v, target=%d\n", nums, target)

	var result [][]int
	var current []int

	var backtrack func(start int, currentSum int)
	backtrack = func(start int, currentSum int) {
		// Если сумма превысила цель — отсекаем эту ветвь
		if currentSum > target {
			return
		}

		// Если нашли нужную сумму — сохраняем результат
		if currentSum == target {
			temp := make([]int, len(current))
			copy(temp, current)
			result = append(result, temp)
			return
		}

		for i := start; i < len(nums); i++ {
			current = append(current, nums[i])

			backtrack(i, currentSum+nums[i])

			// Откат
			current = current[:len(current)-1]
		}
	}

	backtrack(0, 0)

	fmt.Printf("Выход: %v\n", result)
}
