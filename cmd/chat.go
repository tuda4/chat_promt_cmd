/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat command",
	Long:  `Conversation with GPT-3 using the command line. This is a simple command line tool that allows you to chat with GPT-3 using the command line.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		// set up a channel to listen for interupt signals, to exit
		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-signChan
			fmt.Println("\nSee you next time!")
			os.Exit(0)
		}()

		llm, err := openai.New()
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()

		fmt.Println("Enter initial message:")

		initialMes, _ := reader.ReadString('\n')
		initialMes = strings.TrimSpace(initialMes)

		contents := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, initialMes),
		}
		fmt.Println("Initial message received. Entering chat mode...")
		for {
			fmt.Println("->")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			switch input {
			case "quit", "exit":
				fmt.Println("See you next time!")
				os.Exit(0)
			default:
				response := ""
				contents = append(contents, llms.TextParts(llms.ChatMessageTypeHuman, input))
				llm.GenerateContent(ctx, contents, llms.WithMaxTokens(1024),
					llms.WithStreamingFunc(func(ctx context.Context, chuck []byte) error {
						fmt.Println(string(chuck))
						response += string(chuck)
						return nil
					}),
				)
				contents = append(contents, llms.TextParts(llms.ChatMessageTypeSystem, response))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
