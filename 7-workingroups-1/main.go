package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Ktilis/DiskrMat/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TreeSession хранит дерево и время последнего доступа
type TreeSession struct {
	Tree       *common.TreeNode
	LastAccess time.Time
}

// TreeStore хранит TreeSession по ID сессии
var TreeStore = sync.Map{}

const SessionTimeout = 60 * time.Minute

type SearchRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Target    int    `json:"target"`
}

func main() {
	r := gin.Default()
	rand.Seed(time.Now().UnixNano())

	// Запуск фоновой горутины для очистки
	go cleanupExpiredSessions()

	// Эндпоинты
	api := r.Group("/api")
	{
		api.GET("/tree", handleGetTree)
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

func handleGetTree(c *gin.Context) {
	sessionID := uuid.New().String()
	tree := generateRandomTree(15)
	TreeStore.Store(sessionID, &TreeSession{
		Tree:       tree,
		LastAccess: time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"sessionId": sessionID,
		"tree":      tree,
	})
}

func handleSearch(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, ok := TreeStore.Load(req.SessionID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Сессия не найдена"})
		return
	}
	session := val.(*TreeSession)
	session.LastAccess = time.Now() // Обновление времени доступа
	TreeStore.Store(req.SessionID, session)

	var path []int
	found := dfs(session.Tree, req.Target, &path)

	c.JSON(http.StatusOK, gin.H{
		"path":  path,
		"found": found,
	})
}

func cleanupExpiredSessions() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		now := time.Now()
		TreeStore.Range(func(key, value interface{}) bool {
			session := value.(*TreeSession)
			if now.Sub(session.LastAccess) > SessionTimeout {
				TreeStore.Delete(key)
			}
			return true
		})
	}
}

func generateRandomTree(n int) *common.TreeNode {
	if n <= 0 {
		return nil
	}

	used := make(map[int]bool)
	nodes := make([]*common.TreeNode, n)
	for i := 0; i < n; i++ {
		val := rand.Intn(100)
		for used[val] {
			val = rand.Intn(100)
		}
		used[val] = true
		nodes[i] = &common.TreeNode{Val: val}
	}

	// Случайное соединение узлов для формирования бинарного дерева
	for i := 1; i < n; i++ {
		for {
			parentIdx := rand.Intn(i)
			if nodes[parentIdx].Left == nil {
				nodes[parentIdx].Left = nodes[i]
				break
			} else if nodes[parentIdx].Right == nil {
				nodes[parentIdx].Right = nodes[i]
				break
			}
		}
	}

	return nodes[0]
}

func dfs(node *common.TreeNode, target int, path *[]int) bool {
	if node == nil {
		return false
	}

	*path = append(*path, node.Val)

	if node.Val == target {
		return true
	}

	if dfs(node.Left, target, path) {
		return true
	}

	if dfs(node.Right, target, path) {
		return true
	}

	return false
}
