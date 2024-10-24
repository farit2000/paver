package paver

import (
	"fmt"
	"os/exec"
)

type Manifest struct {
	PackageName string   `yaml:"package"`
	Deps        []string `yaml:"deps"`
	Scripts     []string `yaml:"scripts"`
}

func (p *Manifest) GetID() string {
	return p.PackageName
}

func (p *Manifest) Run() (interface{}, error) {
	fmt.Println("Running package:", p.PackageName)
	outputs := make([][]byte, 0, len(p.Scripts))
	for _, script := range p.Scripts {
		fmt.Println("Running script:", script)
		out, err := exec.Command("sh", script).Output()
		if err != nil {
			return out, fmt.Errorf("failed to run script %s: %v", script, err)
		}
		outputs = append(outputs, out)
	}
	return outputs, nil
}

type edge struct {
	fromPackage, toPackage string
}

func (e *edge) GetFromID() string {
	return e.fromPackage
}

func (e *edge) GetToID() string {
	return e.toPackage
}
