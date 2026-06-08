package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Ktilis/DiskrMatCollege/common"
)

type Point struct {
	row, col int
}

func main() {
	fmt.Println("DiskrMat: Поиск в ширину")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nКакую задачу выполнить?")
	fmt.Println("1 - Задача 1: Кратчайший путь в лабиринте (по матрице)")
	fmt.Println("2 - Задача 2: Минимальные ходы коня")
	fmt.Println("3 - Задача 3: Гнилые апельсины")
	fmt.Println("4 - Задача 4: Преобразование строки")
	fmt.Println("5 - Задача 5: Ближайший выход из лабиринта")
	fmt.Println("6 - Задача 6: Поиск слова в сетке")
	fmt.Println("7 - Задача 7: Острова (Количество компонент связности)")
	fmt.Println("8 - Задача 8: Минимальный мост между островами")
	fmt.Println("9 - Задача 9: Вращающиеся замки")
	fmt.Println("10 - Задача 10: Эвакуация из комнат")
	fmt.Println("11 - Задача 11: Снежный ком")
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
		solveTask6()
	case "7":
		solveTask7()
	case "8":
		solveTask8()
	case "9":
		solveTask9()
	case "10":
		solveTask10()
	case "11":
		solveTask11()
	default:
		fmt.Println("Неизвестная задача.")
	}
}

// ==========================================
// Задача 1: Кратчайший путь в лабиринте
// ==========================================

func solveTask1() {
	n := 3
	maze := [][]int{
		{0, 0, 0},
		{1, 1, 0},
		{0, 0, 0},
	}
	start := Point{0, 0}
	end := Point{2, 2}

	fmt.Printf("Размер лабиринта: %dx%d\n", n, n)
	fmt.Printf("Старт: (%d, %d), Финиш: (%d, %d)\n", start.row, start.col, end.row, end.col)

	steps := bfsMaze(maze, start, end)
	fmt.Printf("Выход: %d\n", steps)
}

func bfsMaze(maze [][]int, start Point, end Point) int {
	n := len(maze)
	if n == 0 {
		return -1
	}
	m := len(maze[0])

	dirs := []Point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	queue := []Point{start}

	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	dist[start.row][start.col] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.row == end.row && curr.col == end.col {
			return dist[curr.row][curr.col]
		}

		for _, d := range dirs {
			nextRow, nextCol := curr.row+d.row, curr.col+d.col

			if nextRow >= 0 && nextRow < n && nextCol >= 0 && nextCol < m {
				if maze[nextRow][nextCol] == 0 && dist[nextRow][nextCol] == -1 {
					dist[nextRow][nextCol] = dist[curr.row][curr.col] + 1
					queue = append(queue, Point{nextRow, nextCol})
				}
			}
		}
	}

	return -1
}

// ==========================================
// Задача 2: Минимальные ходы коня
// ==========================================

func solveTask2() {
	start := Point{0, 0}
	end := Point{7, 7}

	fmt.Println("Шахматная доска: 8x8")
	fmt.Printf("Старт коня: (%d, %d), Финиш: (%d, %d)\n", start.row, start.col, end.row, end.col)

	steps := bfsKnight(start, end)
	fmt.Printf("Выход: %d\n", steps)
}

func bfsKnight(start Point, end Point) int {
	dirs := []Point{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
		{1, -2}, {1, 2}, {2, -1}, {2, 1},
	}

	queue := []Point{start}

	dist := make([][]int, 8)
	for i := range dist {
		dist[i] = make([]int, 8)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	dist[start.row][start.col] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.row == end.row && curr.col == end.col {
			return dist[curr.row][curr.col]
		}

		for _, d := range dirs {
			nextRow, nextCol := curr.row+d.row, curr.col+d.col

			if nextRow >= 0 && nextRow < 8 && nextCol >= 0 && nextCol < 8 {
				if dist[nextRow][nextCol] == -1 {
					dist[nextRow][nextCol] = dist[curr.row][curr.col] + 1
					queue = append(queue, Point{nextRow, nextCol})
				}
			}
		}
	}

	return -1
}

// ==========================================
// Задача 3: Гнилые апельсины
// ==========================================

func solveTask3() {
	grid := [][]int{
		{2, 1, 1},
		{1, 1, 0},
		{0, 1, 1},
	}

	fmt.Println("Матрица апельсинов:")
	for _, row := range grid {
		fmt.Println(row)
	}

	minutes := orangesRotting(grid)
	fmt.Printf("Выход: %d\n", minutes)
}

func orangesRotting(grid [][]int) int {
	rows := len(grid)
	if rows == 0 {
		return -1
	}
	cols := len(grid[0])

	queue := []Point{}
	freshCount := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 2 {
				queue = append(queue, Point{r, c})
			} else if grid[r][c] == 1 {
				freshCount++
			}
		}
	}

	if freshCount == 0 {
		return 0
	}

	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	minutes := 0

	for len(queue) > 0 && freshCount > 0 {
		minutes++
		levelSize := len(queue)

		for i := 0; i < levelSize; i++ {
			curr := queue[0]
			queue = queue[1:]

			for _, d := range dirs {
				nextRow, nextCol := curr.row+d.row, curr.col+d.col

				if nextRow >= 0 && nextRow < rows && nextCol >= 0 && nextCol < cols {
					if grid[nextRow][nextCol] == 1 {
						grid[nextRow][nextCol] = 2
						freshCount--
						queue = append(queue, Point{nextRow, nextCol})
					}
				}
			}
		}
	}

	if freshCount > 0 {
		return -1
	}

	return minutes
}

// ==========================================
// Задача 4: Преобразование строки
// ==========================================

type wordNode struct {
	word  string
	steps int
}

func solveTask4() {
	A := "hit"
	B := "cog"
	words := []string{"hot", "dot", "dog", "lot", "log", "cog"}

	fmt.Printf("Начальная строка A: %s\n", A)
	fmt.Printf("Конечная строка B:  %s\n", B)
	fmt.Printf("Список слов: %v\n", words)

	steps := wordLadder(A, B, words)
	fmt.Printf("Выход: %d\n", steps)
}

func wordLadder(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}

	if !wordSet[endWord] {
		return -1
	}

	queue := []wordNode{{word: beginWord, steps: 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.word == endWord {
			return curr.steps
		}

		wordBytes := []byte(curr.word)
		for i := 0; i < len(wordBytes); i++ {
			originalChar := wordBytes[i]

			for c := byte('a'); c <= byte('z'); c++ {
				if c == originalChar {
					continue
				}

				wordBytes[i] = c
				newWord := string(wordBytes)

				if wordSet[newWord] {
					queue = append(queue, wordNode{word: newWord, steps: curr.steps + 1})
					delete(wordSet, newWord)
				}
			}
			wordBytes[i] = originalChar
		}
	}

	return -1
}

// ==========================================
// Задача 5: Ближайший выход из лабиринта
// ==========================================

func solveTask5() {
	maze := [][]int{
		{1, 0, 0, 1},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
	}
	start := Point{1, 1}

	fmt.Println("Матрица лабиринта:")
	for _, row := range maze {
		fmt.Println(row)
	}
	fmt.Printf("Старт: (%d, %d)\n", start.row, start.col)

	steps := nearestExit(maze, start)
	fmt.Printf("Выход: %d\n", steps)
}

func nearestExit(maze [][]int, start Point) int {
	rows := len(maze)
	if rows == 0 {
		return -1
	}
	cols := len(maze[0])

	dirs := []Point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	queue := []Point{start}

	dist := make([][]int, rows)
	for i := range dist {
		dist[i] = make([]int, cols)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	dist[start.row][start.col] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		isAtBorder := curr.row == 0 || curr.row == rows-1 || curr.col == 0 || curr.col == cols-1
		if isAtBorder && (curr.row != start.row || curr.col != start.col) {
			return dist[curr.row][curr.col]
		}

		for _, d := range dirs {
			nextRow, nextCol := curr.row+d.row, curr.col+d.col

			if nextRow >= 0 && nextRow < rows && nextCol >= 0 && nextCol < cols {
				if maze[nextRow][nextCol] == 0 && dist[nextRow][nextCol] == -1 {
					dist[nextRow][nextCol] = dist[curr.row][curr.col] + 1
					queue = append(queue, Point{nextRow, nextCol})
				}
			}
		}
	}

	return -1
}

// ==========================================
// Задача 6: Поиск слова в сетке
// ==========================================

type wordSearchState struct {
	r, c  int
	index int
	vis   map[Point]bool
}

func solveTask6() {
	grid := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	word := "ABCCED"

	fmt.Println("Сетка букв:")
	for _, row := range grid {
		fmt.Printf("  ")
		for _, char := range row {
			fmt.Printf("'%c' ", char)
		}
		fmt.Println()
	}
	fmt.Printf("Искомое слово: %s\n", word)

	found := existWordBFS(grid, word)
	fmt.Printf("Выход: %t\n", found)
}

func existWordBFS(board [][]byte, word string) bool {
	rows := len(board)
	if rows == 0 || len(word) == 0 {
		return false
	}
	cols := len(board[0])

	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queue := []wordSearchState{}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == word[0] {
				vis := make(map[Point]bool)
				vis[Point{i, j}] = true
				queue = append(queue, wordSearchState{r: i, c: j, index: 0, vis: vis})
			}
		}
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.index == len(word)-1 {
			return true
		}

		for _, d := range dirs {
			nextRow, nextCol := curr.r+d.row, curr.c+d.col
			nextPoint := Point{nextRow, nextCol}

			if nextRow >= 0 && nextRow < rows && nextCol >= 0 && nextCol < cols {
				if board[nextRow][nextCol] == word[curr.index+1] && !curr.vis[nextPoint] {
					newVis := make(map[Point]bool)
					for k, v := range curr.vis {
						newVis[k] = v
					}
					newVis[nextPoint] = true

					queue = append(queue, wordSearchState{
						r:     nextRow,
						c:     nextCol,
						index: curr.index + 1,
						vis:   newVis,
					})
				}
			}
		}
	}

	return false
}

// ==========================================
// Задача 7: Острова (Количество компонент связности)
// ==========================================

func solveTask7() {
	grid := [][]int{
		{1, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
	}

	fmt.Println("Карта (1 - земля, 0 - вода):")
	for _, row := range grid {
		fmt.Println(row)
	}

	count := numIslands(grid)
	fmt.Printf("Выход: %d (количество островов)\n", count)
}

func numIslands(grid [][]int) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	islands := 0
	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 1 {
				islands++
				queue := []Point{{r, c}}
				grid[r][c] = 0

				for len(queue) > 0 {
					curr := queue[0]
					queue = queue[1:]

					for _, d := range dirs {
						nr, nc := curr.row+d.row, curr.col+d.col
						if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == 1 {
							grid[nr][nc] = 0
							queue = append(queue, Point{nr, nc})
						}
					}
				}
			}
		}
	}
	return islands
}

// ==========================================
// Задача 8: Минимальный мост между островами
// ==========================================

func solveTask8() {
	grid := [][]int{
		{1, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
	}

	fmt.Println("Карта островов:")
	for _, row := range grid {
		fmt.Println(row)
	}

	bridge := shortestBridge(grid)
	fmt.Printf("Выход: %d (длина моста)\n", bridge)
}

func shortestBridge(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queue := []Point{}
	found := false

	for r := 0; r < rows && !found; r++ {
		for c := 0; c < cols && !found; c++ {
			if grid[r][c] == 1 {
				dfsBridge(grid, r, c, &queue)
				found = true
			}
		}
	}

	steps := 0
	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			curr := queue[0]
			queue = queue[1:]

			for _, d := range dirs {
				nr, nc := curr.row+d.row, curr.col+d.col
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
					if grid[nr][nc] == 1 {
						return steps
					}
					if grid[nr][nc] == 0 {
						grid[nr][nc] = 2
						queue = append(queue, Point{nr, nc})
					}
				}
			}
		}
		steps++
	}
	return -1
}

func dfsBridge(grid [][]int, r, c int, queue *[]Point) {
	if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] != 1 {
		return
	}
	grid[r][c] = 2
	*queue = append(*queue, Point{r, c})
	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range dirs {
		dfsBridge(grid, r+d.row, c+d.col, queue)
	}
}

// ==========================================
// Задача 9: Вращающиеся замки
// ==========================================

type lockState struct {
	code  string
	steps int
}

func solveTask9() {
	start := "0000"
	target := "0202"
	forbidden := []string{"0001", "0101"}

	fmt.Printf("Старт: %s, Цель: %s\nЗапрещены: %v\n", start, target, forbidden)

	steps := openLock(start, target, forbidden)
	fmt.Printf("Выход: %d\n", steps)
}

func openLock(start string, target string, deadends []string) int {
	deadMap := make(map[string]bool)
	for _, d := range deadends {
		deadMap[d] = true
	}

	if deadMap[start] {
		return -1
	}

	queue := []lockState{{code: start, steps: 0}}
	visited := make(map[string]bool)
	visited[start] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.code == target {
			return curr.steps
		}

		for i := 0; i < 4; i++ {
			for _, dir := range []int{-1, 1} {
				runes := []rune(curr.code)
				digit := int(runes[i] - '0')
				nextDigit := (digit + dir + 10) % 10
				runes[i] = rune(nextDigit + '0')
				nextCode := string(runes)

				if !visited[nextCode] && !deadMap[nextCode] {
					visited[nextCode] = true
					queue = append(queue, lockState{code: nextCode, steps: curr.steps + 1})
				}
			}
		}
	}
	return -1
}

// ==========================================
// Задача 10: Эвакуация из комнат
// ==========================================

func solveTask10() {
	grid := [][]string{
		{"S", "0", "1", "E"},
		{"1", "0", "1", "0"},
		{"1", "0", "0", "S"},
	}

	fmt.Println("План здания:")
	for _, row := range grid {
		fmt.Println(row)
	}

	time := evacuate(grid)
	fmt.Printf("Выход: %d (минимальное время для всех)\n", time)
}

func evacuate(grid [][]string) int {
	rows, cols := len(grid), 0
	if rows > 0 {
		cols = len(grid[0])
	}

	queue := []Point{}
	dist := make([][]int, rows)
	totalStarts := 0

	for r := 0; r < rows; r++ {
		dist[r] = make([]int, cols)
		for c := 0; c < cols; c++ {
			dist[r][c] = -1
			if grid[r][c] == "E" {
				queue = append(queue, Point{r, c})
				dist[r][c] = 0
			} else if grid[r][c] == "S" {
				totalStarts++
			}
		}
	}

	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	maxTime := 0
	foundStarts := 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if grid[curr.row][curr.col] == "S" {
			foundStarts++
			if dist[curr.row][curr.col] > maxTime {
				maxTime = dist[curr.row][curr.col]
			}
		}

		for _, d := range dirs {
			nr, nc := curr.row+d.row, curr.col+d.col
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
				if (grid[nr][nc] == "0" || grid[nr][nc] == "S") && dist[nr][nc] == -1 {
					dist[nr][nc] = dist[curr.row][curr.col] + 1
					queue = append(queue, Point{nr, nc})
				}
			}
		}
	}

	if foundStarts < totalStarts {
		return -1
	}
	return maxTime
}

// ==========================================
// Задача 11: Снежный ком
// ==========================================

type snowballState struct {
	val   int
	steps int
}

func solveTask11() {
	target := 10
	fmt.Printf("Цель N: %d\n", target)

	steps := minSnowballOperations(target)
	fmt.Printf("Выход: %d операций\n", steps)
}

func minSnowballOperations(target int) int {
	if target == 1 {
		return 0
	}

	queue := []snowballState{{val: 1, steps: 0}}
	visited := make(map[int]bool)
	visited[1] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.val == target {
			return curr.steps
		}

		nextVals := []int{curr.val * 2, curr.val * 3, curr.val + 1}

		for _, nextVal := range nextVals {
			if nextVal <= target && !visited[nextVal] {
				visited[nextVal] = true
				queue = append(queue, snowballState{val: nextVal, steps: curr.steps + 1})
			}
		}
	}
	return -1
}
