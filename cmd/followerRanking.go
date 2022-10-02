package cmd

import (
	"fmt"
	"github.com/ibnr2hc/TwitterUtils/cmd/models"
	"github.com/spf13/cobra"
	"os"
)

var followerRankingCmd = &cobra.Command{
	Use:   "followerRanking",
	Short: "指定したユーザーのフォロワーをランキングで表示する",
	Long: `指定したユーザーのフォロワーをランキングで表示する
フォロワーのフォロワー数から上位5名を表示する
もしかしたらあの人は有名人と友達かも`,
	Run: func(cmd *cobra.Command, args []string) {
		// 調査対象のユーザーとそのフォロワー数を出力する。
		screenName := args[0]
		surveyUser := models.NewUser(screenName)
		surveyUser.ShowInfoAsSurveyUser()

		// 鍵垢の場合は処理を終える
		if surveyUser.Protected {
			fmt.Println("[Info] 鍵垢のため調査できません。")
			os.Exit(0)
		}

		surveyUser.Followers = surveyUser.GetFollowers()
		surveyUser.ShowRankingByFollowerCount()
	},
}

func init() {
	rootCmd.AddCommand(followerRankingCmd)
}
