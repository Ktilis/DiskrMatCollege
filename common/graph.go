package common

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Edge struct {
	Node   string
	Weight int
}

type Graph struct {
	Edges map[string][]Edge
}

func NewGraph() *Graph {
	return &Graph{
		Edges: make(map[string][]Edge),
	}
}

func (g *Graph) AddEdge(from, to string, weight int) {
	g.Edges[from] = append(g.Edges[from], Edge{Node: to, Weight: weight})
}

func (g *Graph) AddEdgeBi(from, to string, weight int) {
	g.Edges[from] = append(g.Edges[from], Edge{Node: to, Weight: weight})
	g.Edges[to] = append(g.Edges[to], Edge{Node: from, Weight: weight})
}

func buildGraphFromReaderTemplate(reader *bufio.Reader, callback func(string, string, int)) {
	inputGraphRegexp := regexp.MustCompile(`([a-zA-Z0-9А-Яа-я]+)\s+([a-zA-Z0-9А-Яа-я]+)\s+([0-9-]+)`)

	fmt.Println(
		"Введите все связи в формате \"Откуда Куда Стоимость(число)\", после каждой связи нажимайте enter.",
		"После последней связи просто нажмите enter без ввода текста",
	)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			break
		}

		if !inputGraphRegexp.MatchString(text) {
			fmt.Println("Это неверная строка")
			continue
		}

		inputMatches := inputGraphRegexp.FindStringSubmatch(text)
		weight, _ := strconv.Atoi(inputMatches[3])

		// Добавляем связь по внешней функции
		callback(inputMatches[1], inputMatches[2], weight)
	}
}

func BuildGraphFromReader(reader *bufio.Reader) *Graph {
	graph := NewGraph()
	buildGraphFromReaderTemplate(reader, graph.AddEdge)
	return graph
}

func BuildBiGraphFromReader(reader *bufio.Reader) *Graph {
	graph := NewGraph()
	buildGraphFromReaderTemplate(reader, graph.AddEdgeBi)
	return graph
}
