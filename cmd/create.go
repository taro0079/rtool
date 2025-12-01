/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Options struct {
	SortNumber   string
	TicketNumber string
	When         string
	Explanation  string
}

var (
	o = &Options{}
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "ddlファイルを作成します。",
	Long: `リピストのddlファイル作成ルールに基づいて空のddlファイルを作成します。
		ticket_number, とwhenをオプションとして与えてください。
		`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := fmt.Sprintf("%s_%s_%s_%s_%s.sql", time.Now().Format("20060102"), o.SortNumber, o.TicketNumber, o.When, o.Explanation)
		_, err := os.Create(filename)

		if err != nil {
			fmt.Println("書き込みエラー:", err)
		}

		fmt.Println("ファイルが作成されました: ", filename)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&o.SortNumber, "sort_number", "s", "0000", "複数のsqlファイルがあり、適用順序に制約がある場合に指定する。数値の小さい順番にsqlが実行される。")
	createCmd.Flags().StringVarP(&o.TicketNumber, "ticket_number", "t", "", "redmineのチケット番号です")
	createCmd.Flags().StringVarP(&o.When, "when", "w", "", "いつにsqlを実行するかを指定します")
	createCmd.Flags().StringVarP(&o.Explanation, "explanation", "e", "", "sqlの内容の簡単な説明。小文字英数字とアンダースコアのみ。create, update, dropといったsqlの操作とテーブル名を連結する形")
	createCmd.MarkFlagRequired("ticket_number")
	createCmd.MarkFlagRequired("when")
	createCmd.MarkFlagRequired("explanation")
}
