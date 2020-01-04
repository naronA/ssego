package commands

import (
	"fmt"
	"ssego"
	"strings"

	"github.com/urfave/cli"
)

// 検索を実行するコマンド
var searchCommand = cli.Command{
	Name:      "search",
	Usage:     "search documents",
	ArgsUsage: `<query>`,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "number, n",
			Value: 10,
		},
	},
	Action: search,
}

func search(c *cli.Context) error {
	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}
	query := c.Args().Get(0)
	result, err := engine.Search(query, c.Int("number"), "TFIDF")
	if err != nil {
		return err
	}
	printResult(result)
	return nil
}

// 検索結果を表示する
func printResult(results []*ssego.SearchResult) {
	if len(results) == 0 {
		fmt.Println("0 match!!")
		return
	}
	s := make([]string, len(results))
	for i, result := range results {
		s[i] = fmt.Sprintf("rank: %3d, score: %4f, title: %s", i+1, result.Score, result.Title)
	}
	fmt.Println(strings.Join(s, "\n"))
}
