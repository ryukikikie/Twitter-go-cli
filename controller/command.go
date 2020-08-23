package controller

import (
	"fmt"

	"github.com/ryukikikie/twitter-go-cli/models"
)

func Help() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("        <command> [arguments]")
	fmt.Println()
	fmt.Println("The commands are:")
	fmt.Println()
	for name, command := range models.Commands {
		fmt.Printf("        %s:%s\n", name, command.Description)
		if len(command.Option) > 0 {
			fmt.Println("        options")
		}
		for o, description := range command.Option {
			fmt.Printf("%s : %s\n", o, description)
		}
	}

}
