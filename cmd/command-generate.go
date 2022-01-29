package main

import (
	"github.com/spf13/cobra"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/config"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/idea"
	"github.com/xxxbobrxxx/idea-project-manager/pkg/repository"
)

type GenerateCommand struct {
	config.GlobalFlags
	repository.RepositoryFlags
	idea.Project

	cmd *cobra.Command
}

func NewGenerateCommand() *GenerateCommand {
	command := &GenerateCommand{}

	cmd := &cobra.Command{
		Use:          "generate",
		Aliases:      []string{"gen"},
		Short:        "Generate an IDEA project",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE:         command.Execute,
	}
	command.cmd = cmd

	command.Project.AddFlags(cmd.PersistentFlags())
	command.GlobalFlags.AddFlags(cmd.PersistentFlags())
	command.RepositoryFlags.AddFlags(cmd.PersistentFlags())

	_ = command.cmd.MarkPersistentFlagRequired("config")
	_ = command.cmd.MarkPersistentFlagRequired("idea-sources-root")

	return command
}

func (command *GenerateCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GenerateCommand) Execute(_ *cobra.Command, _ []string) (err error) {
	c, err := command.ReadConfig()
	if err != nil {
		return err
	}

	project := command.Project

	for _, repositoryConfig := range c.RepositoryConfigs {
		r, err := repositoryConfig.NewFromConfig()
		if err != nil {
			return err
		}

		err = r.Init(command.RepositoryFlags)
		if err != nil {
			return err
		}

		_, err = r.Clone()
		if err != nil {
			return err
		}

		project.AddRepository(r)
	}

	err = project.Write()
	if err != nil {
		return err
	}

	return nil
}
