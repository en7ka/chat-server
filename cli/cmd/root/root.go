package root

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd представляет базовую команду при вызове без каких-либо подкоманд
var rootCmd = &cobra.Command{
	Use:   "my-app",
	Short: "Мое cli приложение",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Что-то создает",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Что-то удаляет",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Создаем нового пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("ошибка в получении юзернеймов: %v", err.Error())
		}

		log.Printf("user %s created", usernamesStr)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Удаляет пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("ошибка в получении юзернеймов: %v", err.Error())
		}

		log.Printf("user %s deleted", usernamesStr)
	},
}

// Execute добавляет все дочерние команды к команде root и устанавливает соответствующие флаги.
// Это вызывается с помощью main.main(). Для rootCmd это должно произойти только один раз.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)

	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	if err := createUserCmd.MarkFlagRequired("username"); err != nil {
		log.Fatalf("не удалось пометить флаг пользователя как требуется: %s\n", err.Error())
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	if err := deleteUserCmd.MarkFlagRequired("username"); err != nil {
		log.Fatalf("не удалось пометить флаг пользователя как требуется: %s\n", err.Error())
	}
}
