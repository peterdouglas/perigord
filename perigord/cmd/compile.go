// Copyright © 2017 Swarm Market <info@swarm.market>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile contract source files",
	Run: func(cmd *cobra.Command, args []string) {
		err := RunInRoot(func() error {
			if err := compileContracts(); err != nil {
				return err
			}

			return generateBindings()
		})
		if err != nil {
			Fatal(err)
		}
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "(alias for compile)",
	Run:   compileCmd.Run,
}

func init() {
	RootCmd.AddCommand(compileCmd)
	RootCmd.AddCommand(buildCmd)
}

func compileContracts() error {
	// TODO: Figure out relative imports and if we need to do anything else here
	matches, err := filepath.Glob(ContractsDirectory + "/*.sol")
	if err != nil {
		return err
	}

	for _, match := range matches {
		compileContract(match)
	}

	return nil
}

func compileContract(path string) error {
	// TODO: This just shells out atm, could directly integrate abigen and call
	// into it as a library later
	command := "solc"
	args := []string{path, "--bin", "--abi", "--optimize", "--overwrite", "-o", BuildDirectory}
	return ExecWithOutput(command, args...)
}

func generateBindings() error {
	matches, err := filepath.Glob(BuildDirectory + "/*.abi")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(BindingsDirectory, os.FileMode(0755)); err != nil {
		return err
	}

	for _, match := range matches {
		generateBinding(strings.TrimSuffix(match, filepath.Ext(match)))
	}

	return nil
}

func generateBinding(path string) error {
	// TODO: Allow alternate binding directories / package names, in config file
	command := "abigen"
	name := filepath.Base(path)
	abifile := path + ".abi"
	binfile := path + ".bin"
	outfile := filepath.Join(BindingsDirectory, filepath.Base(name)) + ".go"
	args := []string{"--abi", abifile, "--bin", binfile, "--pkg", "bindings", "--type", name, "--out", outfile}
	return ExecWithOutput(command, args...)
}
