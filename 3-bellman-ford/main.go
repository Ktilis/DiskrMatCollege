package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/Ktilis/DiskrMat/common"
)

func main() {
	fmt.Println("DiskrMat: Алгоритм Беллмана-Форда")

	reader := bufio.NewReader(os.Stdin)
	graph := common.BuildGraphFromReader(reader)

	fmt.Print("\nВведите стартовую вершину: ")
	startPoint := common.ReadLine(reader)

	fmt.Println("\nКакую задачу выполнить?")
	fmt.Println("1 - Задача 1: Базовый пример с отрицательным весом")
	fmt.Println("2 - Задача 2: Обнаружение цикла отрицательного веса")
	fmt.Println("3 - Задача 3: Кратчайший путь с отрицательными рёбрами")
	fmt.Println("5 - Задача 5: Восстановление пути")
	fmt.Println("6 - Задача 6: Недостижимый отрицательный цикл")
	fmt.Print("Ваш выбор: ")
	taskChoice := common.ReadLine(reader)

	var endPoint string
	// Конечная вершина нужна для 3 и 5 задачи
	if taskChoice == "3" || taskChoice == "5" {
		fmt.Print("Введите конечную вершину: ")
		endPoint = common.ReadLine(reader)
	}

	fmt.Println("-------------------------")
	if taskChoice == "1" {
		solveTask1(graph, startPoint)
	} else if taskChoice == "2" {
		solveTask2(graph, startPoint)
	} else if taskChoice == "3" || taskChoice == "5" {
		solveTask3(graph, startPoint, endPoint)
	} else if taskChoice == "6" {
		solveTask6(graph, startPoint)
	} else {
		fmt.Println("Неизвестная задача.")
	}
}

func getVertices(g *common.Graph) []string {
	vMap := make(map[string]bool)
	for u, edges := range g.Edges {
		vMap[u] = true
		for _, e := range edges {
			vMap[e.Node] = true
		}
	}
	var vertices []string
	for v := range vMap {
		vertices = append(vertices, v)
	}
	sort.Strings(vertices)
	return vertices
}

func runBellmanFord(g *common.Graph, start string) (map[string]int, map[string]string, bool) {
	distances := make(map[string]int)
	previous := make(map[string]string)
	vertices := getVertices(g)

	for _, v := range vertices {
		distances[v] = math.MaxInt32
	}

	distances[start] = 0
	V := len(vertices)

	for i := 0; i < V-1; i++ {
		updated := false
		for u, edges := range g.Edges {
			if distances[u] == math.MaxInt32 {
				continue
			}
			for _, edge := range edges {
				if distances[u]+edge.Weight < distances[edge.Node] {
					distances[edge.Node] = distances[u] + edge.Weight
					previous[edge.Node] = u
					updated = true
				}
			}
		}
		if !updated {
			break
		}
	}

	hasNegativeCycle := false
	for u, edges := range g.Edges {
		if distances[u] == math.MaxInt32 {
			continue
		}
		for _, edge := range edges {
			if distances[u]+edge.Weight < distances[edge.Node] {
				hasNegativeCycle = true
				break
			}
		}
		if hasNegativeCycle {
			break
		}
	}

	return distances, previous, hasNegativeCycle
}

func getWeight(g *common.Graph, from, to string) int {
	for _, e := range g.Edges[from] {
		if e.Node == to {
			return e.Weight
		}
	}
	return 0
}

func reconstructPath(start, end string, previous map[string]string) []string {
	var path []string
	curr := end
	for curr != "" {
		path = append([]string{curr}, path...)
		if curr == start {
			break
		}
		curr = previous[curr]
	}
	return path
}

func buildFormula(g *common.Graph, path []string, distance int) string {
	var terms []string
	for i := 0; i < len(path)-1; i++ {
		w := getWeight(g, path[i], path[i+1])
		if w < 0 {
			terms = append(terms, fmt.Sprintf("(%d)", w))
		} else {
			terms = append(terms, fmt.Sprintf("%d", w))
		}
	}
	return fmt.Sprintf("%s = %d", strings.Join(terms, " + "), distance)
}

func solveTask1(g *common.Graph, start string) {
	distances, previous, hasCycle := runBellmanFord(g, start)

	if hasCycle {
		fmt.Println("Внимание: в графе обнаружен цикл отрицательного веса! Найти кратчайшие пути невозможно.")
		return
	}

	fmt.Println("Выходные данные:")
	for _, v := range getVertices(g) {
		if distances[v] == math.MaxInt32 {
			fmt.Printf("%s: недостижима\n", v)
			continue
		}

		if v == start {
			fmt.Printf("%s: 0\n", v)
			continue
		}

		path := reconstructPath(start, v, previous)

		if len(path) == 2 {
			fmt.Printf("%s: %d\n", v, distances[v])
		} else {
			pathStr := strings.Join(path, " -> ")
			formulaStr := buildFormula(g, path, distances[v])
			fmt.Printf("%s: %d  (%s, %s)\n", v, distances[v], pathStr, formulaStr)
		}
	}
}

func solveTask2(g *common.Graph, start string) {
	_, _, hasCycle := runBellmanFord(g, start)

	cycleResult := "False"
	if hasCycle {
		cycleResult = "True"
	}

	fmt.Println("Выходные данные:")
	fmt.Printf("Обнаружен цикл отрицательного веса: %s\n", cycleResult)
}

func solveTask3(g *common.Graph, start, end string) {
	distances, previous, hasCycle := runBellmanFord(g, start)

	if hasCycle {
		fmt.Println("Внимание: в графе обнаружен цикл отрицательного веса! Найти кратчайшие пути невозможно.")
		return
	}

	fmt.Println("Выходные данные:")

	if distances[end] == math.MaxInt32 {
		fmt.Printf("Путь от %s до %s не найден.\n", start, end)
		return
	}

	path := reconstructPath(start, end, previous)
	pathStr := strings.Join(path, " -> ")

	fmt.Printf("Длина: %d\n", distances[end])
	fmt.Printf("Путь: %s\n", pathStr)

	if len(path) > 1 {
		formulaStr := buildFormula(g, path, distances[end])
		fmt.Printf("(%s)\n", formulaStr)
	}
}

func solveTask6(g *common.Graph, start string) {
	distances, _, hasCycle := runBellmanFord(g, start)

	// Алгоритм сообщит о цикле только если он ДОСТИЖИМ из стартовой точки.
	// Для недостижимых циклов он проигнорирует их и пойдет дальше.
	if hasCycle {
		fmt.Println("Внимание: в достижимой части графа обнаружен цикл отрицательного веса!")
		return
	}

	fmt.Println("Выходные данные:")
	for _, v := range getVertices(g) {
		if distances[v] == math.MaxInt32 {
			// Выводим знак бесконечности для недостижимых вершин
			fmt.Printf("%s : ∞\n", v)
		} else {
			fmt.Printf("%s : %d\n", v, distances[v])
		}
	}
}
