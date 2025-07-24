package ports

type CommandHandler interface {
	Process(input string) string
}