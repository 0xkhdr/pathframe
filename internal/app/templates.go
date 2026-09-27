package app

import "github.com/0xkhdr/pathframe/internal/artifacts"

func (s Service) Template(mode artifacts.Mode, kind string) (artifacts.Instructions, string, error) {
	instructions, err := artifacts.ArtifactInstructions(mode, kind)
	if err != nil {
		return artifacts.Instructions{}, "", err
	}
	template, err := artifacts.Template(mode, kind)
	return instructions, template, err
}
