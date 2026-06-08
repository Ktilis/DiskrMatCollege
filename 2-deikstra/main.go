package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/Ktilis/DiskrMatCollege/common"
)

func main() {
	fmt.Println("DiskrMat: Алгоритм Дейкстры")

	reader := bufio.NewReader(os.Stdin)

	graph := common.BuildBiGraphFromReader(reader)

	fmt.Println("\nКакую задачу выполнить?")
	fmt.Println("1 - Задача 1. Восстановление пути (один путь)")
	fmt.Println("2 - Задача 2. Граф с несколькими путями одинаковой длины")
	fmt.Println("3 - Задача 3. Проверка на недостижимость")
	fmt.Print("Ваш выбор (1, 2 или 3): ")
	taskChoice := common.ReadLine(reader)

	fmt.Print("Теперь имя начального графа (вершины): ")
	startPoint := common.ReadLine(reader)

	var endPoint string
	// Конечная вершина нужна только для первых двух задач
	if taskChoice == "1" || taskChoice == "2" {
		fmt.Print("Теперь имя конечного графа (вершины): ")
		endPoint = common.ReadLine(reader)
	}

	fmt.Println("-------------------------")
	if taskChoice == "1" {
		dist, path := deikstraTask1(graph, startPoint, endPoint)
		if dist == -1 {
			fmt.Printf("Путь от %s до %s не найден.\n", startPoint, endPoint)
		} else {
			fmt.Printf("Длина: %d\n", dist)
			fmt.Printf("Путь: %s\n", path)
		}
	} else if taskChoice == "2" {
		dist, paths := deikstraTask2(graph, startPoint, endPoint)
		if dist == -1 || len(paths) == 0 {
			fmt.Printf("Путь от %s до %s не найден.\n", startPoint, endPoint)
		} else {
			fmt.Printf("Длина: %d\n", dist)
			fmt.Println("Пути:")
			for i, p := range paths {
				fmt.Printf("%d. %s\n", i+1, p)
			}
		}
	} else if taskChoice == "3" {
		distances := deikstraTask3(graph, startPoint)

		// Собираем все вершины в массив и сортируем по алфавиту
		var nodes []string
		for node := range distances {
			nodes = append(nodes, node)
		}
		sort.Strings(nodes)

		// Выводим результат
		for _, node := range nodes {
			if distances[node] == math.MaxInt32 {
				fmt.Printf("%s: -1\n", node)
			} else {
				fmt.Printf("%s: %d\n", node, distances[node])
			}
		}
	} else {
		fmt.Println("Неизвестная задача.")
	}
}

// deikstraTask1 решает Задачу 1: находит только один кратчайший путь
func deikstraTask1(g *common.Graph, start string, end string) (int, string) {
	distances := make(map[string]int)
	visited := make(map[string]bool)
	previous := make(map[string]string) // Храним только одну предыдущую вершину

	for node := range g.Edges {
		distances[node] = math.MaxInt32
		for _, edge := range g.Edges[node] {
			distances[edge.Node] = math.MaxInt32
		}
	}

	distances[start] = 0

	for i := 0; i < len(distances); i++ {
		currNode := ""
		minDist := math.MaxInt32

		for node, dist := range distances {
			if !visited[node] && dist < minDist {
				minDist = dist
				currNode = node
			}
		}

		if currNode == "" {
			break
		}

		visited[currNode] = true

		for _, edge := range g.Edges[currNode] {
			if !visited[edge.Node] {
				newDist := distances[currNode] + edge.Weight

				if newDist < distances[edge.Node] {
					distances[edge.Node] = newDist
					previous[edge.Node] = currNode
				}
			}
		}
	}

	if distances[end] == math.MaxInt32 {
		return -1, ""
	}

	var path []string
	curr := end

	for curr != "" {
		path = append([]string{curr}, path...)
		if curr == start {
			break
		}
		curr = previous[curr]
	}

	return distances[end], strings.Join(path, " -> ")
}

// deikstraTask2 решает Задачу 2: находит все возможные кратчайшие пути
func deikstraTask2(g *common.Graph, start string, end string) (int, []string) {
	distances := make(map[string]int)
	visited := make(map[string]bool)
	previous := make(map[string][]string) // Храним список вершин для развилок

	for node := range g.Edges {
		distances[node] = math.MaxInt32
		for _, edge := range g.Edges[node] {
			distances[edge.Node] = math.MaxInt32
		}
	}

	distances[start] = 0

	for i := 0; i < len(distances); i++ {
		currNode := ""
		minDist := math.MaxInt32

		for node, dist := range distances {
			if !visited[node] && dist < minDist {
				minDist = dist
				currNode = node
			}
		}

		if currNode == "" {
			break
		}

		visited[currNode] = true

		for _, edge := range g.Edges[currNode] {
			if !visited[edge.Node] {
				newDist := distances[currNode] + edge.Weight

				if newDist < distances[edge.Node] {
					distances[edge.Node] = newDist
					previous[edge.Node] = []string{currNode}
				} else if newDist == distances[edge.Node] {
					previous[edge.Node] = append(previous[edge.Node], currNode)
				}
			}
		}
	}

	if distances[end] == math.MaxInt32 {
		return -1, nil
	}

	allPaths := reconstructPaths(start, end, previous)
	sort.Strings(allPaths)

	return distances[end], allPaths
}

// reconstructPaths рекурсивно восстанавливает все кратчайшие пути от конца к началу
func reconstructPaths(start string, end string, previous map[string][]string) []string {
	var allPaths []string

	// Рекурсивная функция для обхода всех ветвлений (DFS)
	var backtrack func(curr string, currentPath []string)
	backtrack = func(curr string, currentPath []string) {
		// Добавляем текущую вершину в начало пути
		path := append([]string{curr}, currentPath...)

		if curr == start {
			// Если дошли до старта, формируем строку и сохраняем результат
			allPaths = append(allPaths, strings.Join(path, " -> "))
			return
		}

		// Запускаем рекурсию для всех альтернативных предыдущих вершин
		for _, prevNode := range previous[curr] {
			backtrack(prevNode, path)
		}
	}

	// Запускаем восстановление путей начиная с конечной точки
	backtrack(end, []string{})

	return allPaths
}

// deikstraTask3 решает Задачу 3: находит расстояния до всех вершин графа
func deikstraTask3(g *common.Graph, start string) map[string]int {
	distances := make(map[string]int)
	visited := make(map[string]bool)

	// Инициализируем расстояния бесконечностью
	for node := range g.Edges {
		distances[node] = math.MaxInt32
		for _, edge := range g.Edges[node] {
			distances[edge.Node] = math.MaxInt32
		}
	}

	// Если стартовой вершины нет в графе, просто возвращаем пустые/бесконечные расстояния
	if _, exists := distances[start]; exists {
		distances[start] = 0
	}

	for i := 0; i < len(distances); i++ {
		currNode := ""
		minDist := math.MaxInt32

		for node, dist := range distances {
			if !visited[node] && dist < minDist {
				minDist = dist
				currNode = node
			}
		}

		if currNode == "" {
			break
		}

		visited[currNode] = true

		for _, edge := range g.Edges[currNode] {
			if !visited[edge.Node] {
				newDist := distances[currNode] + edge.Weight

				if newDist < distances[edge.Node] {
					distances[edge.Node] = newDist
				}
			}
		}
	}

	return distances
}
