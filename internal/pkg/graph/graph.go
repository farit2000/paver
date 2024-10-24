package graph

import (
	"fmt"
	"sync"
)

type KeyConstraint interface {
	comparable
}

type NodeConstraint[K comparable] interface {
	GetID() K
}

type EdgeConstraint[K comparable] interface {
	GetFromID() K
	GetToID() K
}

type GraphWorker[N NodeConstraint[K], E EdgeConstraint[K], K KeyConstraint] interface {
	Load(nodes []N, edges []E)
	Reset()
	Validate() error

	GetPendingNode() (result N, ok bool)
	ReleaseNode(nodeID K)
}

type Graph[N NodeConstraint[K], E EdgeConstraint[K], K KeyConstraint] struct {
	rwLock sync.RWMutex

	nodes map[K]*Node[N, K]
	edges map[K][]*Edge[N, E, K]

	roots       []*Node[N, K]
	singleNodes map[K]struct{}

	totalNodes uint32

	statusLock   sync.Mutex
	pendingNodes map[K]struct{}
	runningNodes map[K]struct{}
	awaitNodes   map[K]struct{}
	successNodes map[K]struct{}
}

type Node[N NodeConstraint[K], K KeyConstraint] struct {
	node       *N
	inputNodes []K
}

type Edge[N NodeConstraint[K], E EdgeConstraint[K], K KeyConstraint] struct {
	toNode *Node[N, K]

	edge *E
}

func NewGraph[N NodeConstraint[K], E EdgeConstraint[K], K KeyConstraint]() *Graph[N, E, K] {
	return &Graph[N, E, K]{
		nodes:       make(map[K]*Node[N, K]),
		edges:       make(map[K][]*Edge[N, E, K]),
		singleNodes: make(map[K]struct{}),

		pendingNodes: make(map[K]struct{}),
		runningNodes: make(map[K]struct{}),
		awaitNodes:   make(map[K]struct{}),
		successNodes: make(map[K]struct{}),
	}
}

func (g *Graph[N, E, K]) Load(nodes []N, edges []E) {
	g.rwLock.Lock()
	defer g.rwLock.Unlock()

	rootsMap := make(map[K]struct{})

	for _, node := range nodes {
		g.nodes[node.GetID()] = &Node[N, K]{
			node: &node,
		}
		rootsMap[node.GetID()] = struct{}{}
		g.singleNodes[node.GetID()] = struct{}{}
		g.totalNodes++
	}

	for _, edge := range edges {
		g.edges[edge.GetFromID()] = append(g.edges[edge.GetFromID()], &Edge[N, E, K]{
			toNode: g.nodes[edge.GetToID()],
			edge:   &edge,
		})
		g.nodes[edge.GetToID()].inputNodes = append(g.nodes[edge.GetToID()].inputNodes, edge.GetFromID())
		delete(rootsMap, edge.GetToID())

		delete(g.singleNodes, edge.GetToID())
		delete(g.singleNodes, edge.GetFromID())
	}

	for k := range rootsMap {
		g.roots = append(g.roots, g.nodes[k])
		g.pendingNodes[k] = struct{}{}
	}
}

func (g *Graph[N, E, K]) Reset() {
	g.rwLock.Lock()
	defer g.rwLock.Unlock()

	clear(g.nodes)
	clear(g.edges)
	clear(g.singleNodes)

	clear(g.pendingNodes)
	clear(g.runningNodes)
	clear(g.awaitNodes)
	clear(g.successNodes)
}

func (g *Graph[N, E, K]) Validate() error {
	g.rwLock.RLock()
	defer g.rwLock.RUnlock()
	if g.checkCycle() {
		return fmt.Errorf("cycle detected")
	}
	return nil
}

func (g *Graph[N, E, K]) GetPendingNode() (result N, isLast bool) {
	g.statusLock.Lock()
	defer g.statusLock.Unlock()
	if len(g.pendingNodes) > 0 {
		for id := range g.pendingNodes {
			if node, ok1 := g.nodes[id]; ok1 {
				result = *node.node
				delete(g.pendingNodes, id)
				g.runningNodes[id] = struct{}{}
				if len(g.runningNodes)+len(g.awaitNodes)+len(g.successNodes) == int(g.totalNodes) {
					isLast = true
				}
				return
			}
		}
	}
	return
}

func (g *Graph[N, E, K]) ReleaseNode(nodeID K) {
	g.statusLock.Lock()
	defer g.statusLock.Unlock()
	delete(g.runningNodes, nodeID)
	nextNodes := g.getNextNodesByID(nodeID)
	allowToSetToSuccess := true
	for _, nextNode := range nextNodes {
		nextNodeID := nextNode.GetID()
		previousNodes := g.getAllInputNodesByID(nextNodeID)
		allowToSetToPending := true
		for _, previousNode := range previousNodes {
			if previousNode.GetID() == nodeID {
				continue
			}
			if _, ok := g.successNodes[previousNode.GetID()]; ok {
				continue
			} else if _, ok := g.awaitNodes[previousNode.GetID()]; ok {
				continue
			} else {
				allowToSetToPending = false
				break
			}
		}
		if allowToSetToPending {
			g.pendingNodes[nextNodeID] = struct{}{}
		}
		if _, ok := g.pendingNodes[nextNode.GetID()]; ok {
			continue
		} else if _, ok := g.runningNodes[nextNode.GetID()]; ok {
			continue
		} else {
			allowToSetToSuccess = false
		}
	}
	if !allowToSetToSuccess {
		g.awaitNodes[nodeID] = struct{}{}
	} else {
		g.successNodes[nodeID] = struct{}{}
	}
}

func (g *Graph[N, E, K]) dfs(node K, visited, onStack map[K]struct{}) bool {
	visited[node] = struct{}{}
	onStack[node] = struct{}{}
	for _, edge := range g.edges[node] {
		if _, ok := visited[(*edge.toNode.node).GetID()]; !ok {
			if g.dfs((*edge.toNode.node).GetID(), visited, onStack) {
				return true
			}
		} else if _, ok := onStack[(*edge.toNode.node).GetID()]; ok {
			return true
		}
	}
	delete(onStack, node)
	return false
}

func (g *Graph[N, E, K]) checkCycle() bool {
	g.rwLock.RLock()
	defer g.rwLock.RUnlock()
	visited := make(map[K]struct{})
	onStack := make(map[K]struct{})

	for k := range g.nodes {
		if _, ok := visited[k]; !ok {
			if g.dfs(k, visited, onStack) {
				return true
			}
		}
	}
	return false
}

func (g *Graph[N, E, K]) getAllInputNodesByID(id K) []N {
	g.rwLock.RLock()
	defer g.rwLock.RUnlock()
	if node, ok := g.nodes[id]; ok {
		nodes := make([]N, 0, len(node.inputNodes))
		for _, id := range node.inputNodes {
			if n, ok := g.nodes[id]; ok {
				nodes = append(nodes, *n.node)
			}
		}
		return nodes
	}
	return nil
}

func (g *Graph[N, E, K]) getNextNodesByID(id K) []N {
	g.rwLock.RLock()
	defer g.rwLock.RUnlock()

	if edges, ok := g.edges[id]; ok {
		nodes := make([]N, 0, len(edges))
		for _, edge := range edges {
			nodes = append(nodes, *edge.toNode.node)
		}

		return nodes
	}

	return nil
}
