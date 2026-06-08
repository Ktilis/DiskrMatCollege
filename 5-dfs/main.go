package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Ktilis/DiskrMatCollege/common"
)

func main() {
	fmt.Println("DiskrMat: Поиск в глубину")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nКакую задачу выполнить?")
	fmt.Println("1 - Задача 1: Обход бинарного дерева в глубину (префиксный)")
	fmt.Println("2 - Задача 2: Сумма значений всех узлов")
	fmt.Println("3 - Задача 3: Максимальная глубина дерева")
	fmt.Println("4 - Задача 4: Зеркальное отражение дерева")
	fmt.Println("5 - Задача 5: Поиск элемента в дереве")
	fmt.Println("6 - Задача 13: Поиск выхода в лабиринте")
	fmt.Print("Ваш выбор: ")

	taskChoice := common.ReadLine(reader)

	fmt.Println("-------------------------")
	switch taskChoice {
	case "1":
		solveTask1()
	case "2":
		solveTask2()
	case "3":
		solveTask3()
	case "4":
		solveTask4()
	case "5":
		solveTask5()
	case "6":
		solveTask13()
	default:
		fmt.Println("Неизвестная задача.")
	}
}

// Вспомогательная функция для создания дерева из примеров
func buildTreeTask1() *common.TreeNode {
	root := &common.TreeNode{Val: 1}
	root.Left = &common.TreeNode{Val: 2}
	root.Right = &common.TreeNode{Val: 3}
	root.Left.Left = &common.TreeNode{Val: 4}
	root.Left.Right = &common.TreeNode{Val: 5}
	return root
}

func solveTask1() {
	fmt.Println("Решение Задачи 1: Обход бинарного дерева в глубину (префиксный)")

	// Построение дерева из примера:
	//       1
	//      / \
	//     2   3
	//    / \
	//   4   5
	root := buildTreeTask1()

	var result []int

	// dfs корень -> левый -> правый
	var dfs func(node *common.TreeNode)
	dfs = func(node *common.TreeNode) {
		if node == nil {
			return
		}

		// Обрабатываем корень
		result = append(result, node.Val)

		// Идем в левое поддерево
		dfs(node.Left)

		// Идем в правое поддерево
		dfs(node.Right)
	}

	// Запускаем обход от корня
	dfs(root)

	fmt.Printf("Выход: %v\n", result)
}

func solveTask2() {
	fmt.Println("Решение Задачи 2: Сумма значений всех узлов")
	root := buildTreeTask1()

	var dfs func(node *common.TreeNode) int
	dfs = func(node *common.TreeNode) int {
		if node == nil {
			return 0
		}
		// Возвращаем сумму текущего узла и всех узлов в левом и правом поддеревьях
		return node.Val + dfs(node.Left) + dfs(node.Right)
	}

	sum := dfs(root)
	fmt.Printf("Выход: %d\n", sum)
}

func solveTask3() {
	fmt.Println("Решение Задачи 3: Максимальная глубина дерева")
	root := buildTreeTask1()

	var dfs func(node *common.TreeNode) int
	dfs = func(node *common.TreeNode) int {
		if node == nil {
			return 0
		}

		leftDepth := dfs(node.Left)
		rightDepth := dfs(node.Right)

		// Выбираем максимальную глубину из поддеревьев и прибавляем текущий уровень (+1)
		if leftDepth > rightDepth {
			return leftDepth + 1
		}
		return rightDepth + 1
	}

	depth := dfs(root)
	fmt.Printf("Выход: %d\n", depth)
}

func solveTask4() {
	fmt.Println("Решение Задачи 4: Зеркальное отражение дерева")
	root := buildTreeTask1()

	// Рекурсивная функция для отражения дерева
	var dfs func(node *common.TreeNode)
	dfs = func(node *common.TreeNode) {
		if node == nil {
			return
		}

		// Меняем местами левого и правого потомка
		node.Left, node.Right = node.Right, node.Left

		// Идем вглубь по обновленным веткам
		dfs(node.Left)
		dfs(node.Right)
	}

	dfs(root)

	// Чтобы убедиться, что дерево отразилось, выполним его префиксный обход
	var result []int
	var printDfs func(node *common.TreeNode)
	printDfs = func(node *common.TreeNode) {
		if node == nil {
			return
		}
		result = append(result, node.Val)
		printDfs(node.Left)
		printDfs(node.Right)
	}

	printDfs(root)
	fmt.Printf("Выход (префиксный обход после отражения): %v\n", result)
}

func solveTask5() {
	fmt.Println("Решение Задачи 5: Поиск элемента в дереве")
	root := buildTreeTask1()
	target := 5

	var dfs func(node *common.TreeNode, t int) bool
	dfs = func(node *common.TreeNode, t int) bool {
		if node == nil {
			return false
		}
		// Если нашли целевое значение, возвращаем true
		if node.Val == t {
			return true
		}
		// Ищем в левом ИЛИ в правом поддереве (если найдем хотя бы в одном - вернем true)
		return dfs(node.Left, t) || dfs(node.Right, t)
	}

	found := dfs(root, target)
	fmt.Printf("Вход: target = %d\n", target)
	fmt.Printf("Выход: %t\n", found)
}

func solveTask13() {
	fmt.Println("Решение Задачи 13: Поиск выхода в лабиринте")

	// Данные из примера
	maze := [][]int{
		{0, 1, 0},
		{0, 1, 0},
		{0, 0, 0},
	}
	startRow, startCol := 0, 0
	endRow, endCol := 2, 2

	var dfs func(r, c int) bool
	dfs = func(r, c int) bool {
		// Проверка выхода за границы матрицы
		if r < 0 || r >= len(maze) || c < 0 || c >= len(maze[0]) {
			return false
		}

		// Если наткнулись на стену (1) или уже посещали эту клетку (-1)
		if maze[r][c] == 1 || maze[r][c] == -1 {
			return false
		}

		// Если достигли конечной точки
		if r == endRow && c == endCol {
			return true
		}

		// Помечаем текущую клетку как посещенную
		maze[r][c] = -1

		// Рекурсивно идем в 4 направлениях
		return dfs(r+1, c) || // Вниз
			dfs(r-1, c) || // Вверх
			dfs(r, c+1) || // Вправо
			dfs(r, c-1) // Влево
	}

	result := dfs(startRow, startCol)
	fmt.Printf("Вход: start = (%d, %d), end = (%d, %d)\n", startRow, startCol, endRow, endCol)
	fmt.Printf("Выход: %t\n", result)
}
