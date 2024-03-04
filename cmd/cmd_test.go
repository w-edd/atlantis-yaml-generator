package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/w-edd/atlantis-yaml-generator/pkg/config"
)

func TestInitFlags(t *testing.T) {
	// Create a new root command
	cmd := &cobra.Command{}
	initFlags(cmd)

	// Check that all flags have been initialized
	for _, param := range config.ParameterList {
		flag := cmd.PersistentFlags().Lookup(param.Name)
		if flag == nil {
			t.Errorf("Flag %s was not initialized", param.Name)
		}
	}
}
