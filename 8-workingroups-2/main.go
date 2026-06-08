package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Ktilis/DiskrMatCollege/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GraphSession хранит граф и время последнего доступа
type GraphSession struct {
	Graph      *common.Graph
	LastAccess time.Time
}

var GraphStore = sync.Map{}

const SessionTimeout = 60 * time.Minute

type SearchRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Start     string `json:"start" binding:"required"`
	End       string `json:"end" binding:"required"`
	Algorithm string `json:"algorithm" binding:"required"` // "bfs" или "dfs"
}

func main() {
	r := gin.Default()
	rand.Seed(time.Now().UnixNano())

	// Запуск фоновой очистки сессий
	go cleanupExpiredSessions()

	// API Эндпоинты
	api := r.Group("/api")
	{
		api.GET("/graph", handleGetGraph)
		api.POST("/search", handleSearch)
	}

	// Статические файлы из frontend/dist
	r.StaticFS("/assets", http.Dir("./frontend/dist/assets"))
	r.StaticFile("/", "./frontend/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	r.Run(":8080")
}

func handleGetGraph(c *gin.Context) {
	sessionID := uuid.New().String()
	graph := generateRandomGraph(12, 18) // 12 узлов, 18 ребер
	GraphStore.Store(sessionID, &GraphSession{
		Graph:      graph,
		LastAccess: time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"sessionId": sessionID,
		"graph":     graph.Edges,
	})
}

func handleSearch(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, ok := GraphStore.Load(req.SessionID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Сессия не найдена"})
		return
	}
	session := val.(*GraphSession)
	session.LastAccess = time.Now() // Обновление времени доступа
	GraphStore.Store(req.SessionID, session)

	if req.Algorithm == "bfs" {
		path := bfs(session.Graph, req.Start, req.End)
		c.JSON(http.StatusOK, gin.H{
			"paths": [][]string{path}, // Возвращаем массив путей для единообразия
		})
	} else if req.Algorithm == "dfs" {
		paths := dfsAllPaths(session.Graph, req.Start, req.End)
		c.JSON(http.StatusOK, gin.H{
			"paths": paths,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный алгоритм"})
	}
}

func cleanupExpiredSessions() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		now := time.Now()
		GraphStore.Range(func(key, value interface{}) bool {
			session := value.(*GraphSession)
			if now.Sub(session.LastAccess) > SessionTimeout {
				GraphStore.Delete(key)
			}
			return true
		})
	}
}

func generateRandomGraph(nodeCount, edgeCount int) *common.Graph {
	g := common.NewGraph()
	nodes := make([]string, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = string(rune('A' + i))
		g.Edges[nodes[i]] = []common.Edge{} // Убеждаемся, что узел существует
	}

	existingEdges := make(map[string]bool)

	// Гарантируем связность, сначала создавая простое дерево
	perm := rand.Perm(nodeCount)
	for i := 1; i < nodeCount; i++ {
		u := nodes[perm[i]]
		v := nodes[perm[rand.Intn(i)]]
		g.AddEdgeBi(u, v, 1)
		existingEdges[getEdgeKey(u, v)] = true
	}

	// Добавляем дополнительные случайные ребра
	for i := nodeCount - 1; i < edgeCount; i++ {
		uIdx := rand.Intn(nodeCount)
		vIdx := rand.Intn(nodeCount)
		if uIdx == vIdx {
			continue
		}
		u, v := nodes[uIdx], nodes[vIdx]
		key := getEdgeKey(u, v)
		if !existingEdges[key] {
			g.AddEdgeBi(u, v, 1)
			existingEdges[key] = true
		}
	}

	return g
}

func getEdgeKey(u, v string) string {
	if u < v {
		return u + "-" + v
	}
	return v + "-" + u
}

func bfs(g *common.Graph, start, end string) []string {
	if start == end {
		return []string{start}
	}

	queue := [][]string{{start}}
	visited := map[string]bool{start: true}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		u := path[len(path)-1]

		for _, edge := range g.Edges[u] {
			if edge.Node == end {
				return append(path, edge.Node)
			}
			if !visited[edge.Node] {
				visited[edge.Node] = true
				newPath := make([]string, len(path))
				copy(newPath, path)
				queue = append(queue, append(newPath, edge.Node))
			}
		}
	}

	return nil
}

func dfsAllPaths(g *common.Graph, start, end string) [][]string {
	var paths [][]string
	visited := make(map[string]bool)
	var currentPath []string

	var search func(string)
	search = func(u string) {
		visited[u] = true
		currentPath = append(currentPath, u)

		if u == end {
			pathCopy := make([]string, len(currentPath))
			copy(pathCopy, currentPath)
			paths = append(paths, pathCopy)
		} else {
			for _, edge := range g.Edges[u] {
				if !visited[edge.Node] {
					search(edge.Node)
				}
			}
		}

		// Возврат
		currentPath = currentPath[:len(currentPath)-1]
		visited[u] = false
	}

	search(start)
	return paths
}
