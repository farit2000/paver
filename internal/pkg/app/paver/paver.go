package paver

import (
	"fmt"
	"os"
	"path"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/farit2000/paver/internal/pkg/graph"
)

const manifestFileName = "manifest.yaml"

type paver struct {
	packageDirPath string

	graph graph.GraphWorker[*Manifest, *edge, string]
}

func NewPaver(packageDirPath string) paver {
	return paver{
		packageDirPath: packageDirPath,
		graph:          graph.NewGraph[*Manifest, *edge](),
	}
}

func (p paver) Run(workersNum int) error {
	packages, edges, err := p.dirWalker()
	if err != nil {
		return err
	}
	p.graph.Load(packages, edges)
	if err := p.graph.Validate(); err != nil {
		return fmt.Errorf("graph validation failed: %v", err)
	}

	pool := NewWorkersPool(workersNum)
	resultChan := make(chan TaskResult)
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for result := range resultChan {
			if result.Error != nil {
				fmt.Printf("Error run script: %s, output: %s, err: %v", result.ID, result.Value, result.Error)
				close(stopChan)
				return
			}
			fmt.Printf("Task %s completed with result: %s\n", result.ID, result.Value)
			p.graph.ReleaseNode(result.ID)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				pool.Shutdown()
				close(resultChan)
				return
			default:
			}
			task, isLast := p.graph.GetPendingNode()
			if task == nil {
				continue
			}
			err := pool.Submit(Task{
				ID:       task.GetID(),
				TaskFunc: task.Run,
			}, resultChan)
			if err != nil {
				panic(err)
			}
			if isLast {
				close(stopChan)
			}
		}
	}()
	wg.Wait()

	return nil
}

func (p paver) dirWalker() (packages []*Manifest, edges []*edge, err error) {
	entries, err := os.ReadDir(p.packageDirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error read package dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			var mn Manifest

			yamlFile, err := os.ReadFile(path.Join(p.packageDirPath, e.Name(), manifestFileName))
			if err != nil {
				return nil, nil, fmt.Errorf("read manifest file err: %v ", err)
			}
			err = yaml.Unmarshal(yamlFile, &mn)
			if err != nil {
				return nil, nil, fmt.Errorf("unmarshal: %v", err)
			}

			scriptsFullPaths := make([]string, 0, len(mn.Scripts))
			for _, script := range mn.Scripts {
				scriptsFullPaths = append(scriptsFullPaths, path.Join(p.packageDirPath, e.Name(), script))
			}
			mn.Scripts = scriptsFullPaths

			packages = append(packages, &mn)
			for _, dep := range mn.Deps {
				edges = append(edges, &edge{
					fromPackage: dep,
					toPackage:   mn.PackageName,
				})
			}
		}
	}

	return packages, edges, nil
}
