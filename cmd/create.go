package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

type Options struct {
	SortNumber   string
	TicketNumber string
	When         string
	Explanation  string
	Filepath     string
}

type RequestModelOptions struct {
	Name        string
	Namespace   string
	Mode        string
	WithFactory bool
	Stdout      bool
}

var (
	o  = &Options{}
	ro = &RequestModelOptions{}
)
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "リソースを作成します",
}

var createDdlCmd = &cobra.Command{
	Use:   "ddl",
	Short: "ddlファイルを作成します。",
	Long: `リピストのddlファイル作成ルールに基づいて空のddlファイルを作成します。
		ticket_number, とwhenをオプションとして与えてください。
		`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := fmt.Sprintf("%s_%s_%s_%s_%s.sql", time.Now().Format("20060102"), o.SortNumber, o.TicketNumber, o.When, o.Explanation)
		fullPath := filepath.Join(o.Filepath, filename)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("ディレクトリ作成エラー: %w", err)
		}

		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("書き込みエラー: %w", err)
		}

		defer file.Close()

		fmt.Println("ファイルが作成されました: ", filename)

		return nil
	},
}

var requestModelCmd = &cobra.Command{
	Use:   "requestModel",
	Short: "リクエストモデルのテンプレートを作成します。",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := map[string]string{
			"Namespace": ro.Namespace,
			"Name":      ro.Name,
			"Mode":      ro.Mode,
		}
		t, err := template.New("requestModel.tmpl").ParseFiles("./cmd//templates/requestModel.tmpl")

		if err != nil {
			return fmt.Errorf("テンプレートファイルの読み込みでエラーが発生しました: %w", err)
		}

		var writer io.Writer

		if ro.Stdout {
			writer = os.Stdout
		} else {
			filename := fmt.Sprintf("%s.php", ro.Name)
			f, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("ファイル作成でエラーが発生しました: %w", err)
			}

			defer f.Close()
			writer = f
			fmt.Printf("ファイルを生成しました: %s\n", filename)

		}

		if err = t.Execute(writer, data); err != nil {
			return fmt.Errorf("テンプレートアサインでエラーが発生しました: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createDdlCmd)
	createCmd.AddCommand(requestModelCmd)
	createDdlCmd.Flags().StringVarP(&o.SortNumber, "sort_number", "s", "0000", "複数のsqlファイルがあり、適用順序に制約がある場合に指定する。数値の小さい順番にsqlが実行される。")
	createDdlCmd.Flags().StringVarP(&o.TicketNumber, "ticket_number", "t", "", "redmineのチケット番号です")
	createDdlCmd.Flags().StringVarP(&o.When, "when", "w", "", "いつにsqlを実行するかを指定します")
	createDdlCmd.Flags().StringVarP(&o.Explanation, "explanation", "e", "", "sqlの内容の簡単な説明。小文字英数字とアンダースコアのみ。create, update, dropといったsqlの操作とテーブル名を連結する形")
	createDdlCmd.Flags().StringVarP(&o.Filepath, "file_path", "f", ".", "sqlファイルを出力するファイルパス")
	createDdlCmd.MarkFlagRequired("ticket_number")
	createDdlCmd.MarkFlagRequired("when")
	createDdlCmd.MarkFlagRequired("explanation")

	requestModelCmd.Flags().StringVarP(&ro.Name, "name", "n", "", "リクエストモデルの名前")
	requestModelCmd.MarkFlagRequired("name")
	requestModelCmd.Flags().StringVarP(&ro.Mode, "mode", "m", "default", "リクエストモデルのモード")
	requestModelCmd.Flags().StringVarP(&ro.Namespace, "namespace", "a", "", "リクエストモデルのnamespace. App\\model\\requests\\rpst\\以下を指定してください。")
	requestModelCmd.MarkFlagRequired("namespace")
	requestModelCmd.Flags().BoolVarP(&ro.WithFactory, "with-factory", "f", false, "リクエストモデルファクトリも同時に作成するか")
	requestModelCmd.Flags().BoolVar(&ro.Stdout, "stdout", false, "ファイル作成ではなく標準出力する")
}
