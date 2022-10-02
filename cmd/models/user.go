package models

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type User struct {
	ScreenName    string  // ScreenName
	AccountName   string  // アカウント名
	Protected     bool    // 鍵垢か
	FollowerCount int     // フォロワー数
	Followers     []*User // フォロワー一覧
}

// TODO: 命名を考え直したい。initUser?
func NewUser(screenName string) *User {
	client := getTwitterClient()

	// ユーザー情報の取得
	user, _, err := client.Users.Show(&twitter.UserShowParams{ScreenName: screenName})
	if err != nil {
		fmt.Println("[Error]: " + err.Error())
		os.Exit(1)
	}

	return &User{
		ScreenName:    screenName,
		AccountName:   user.Name,
		FollowerCount: user.FollowersCount,
		Protected:     user.Protected,
	}
}

// ShowRankingByFollowerCount
// フォロワー数のランキングを表示する
// ---------------------------
func (u *User) ShowRankingByFollowerCount() {
	sort.Slice(u.Followers,
		func(i, j int) bool { return u.Followers[i].FollowerCount > u.Followers[j].FollowerCount })

	// TODO: forで回す
	// TODO: 5人以上のフォロワーがいない場合にエラーとしない
	u.Followers[0].ShowInfoAsFollowerUser(1)
	u.Followers[1].ShowInfoAsFollowerUser(2)
	u.Followers[2].ShowInfoAsFollowerUser(3)
	u.Followers[3].ShowInfoAsFollowerUser(4)
	u.Followers[4].ShowInfoAsFollowerUser(5)
}

// GetFollowers
// フォロワーのリストを返す
// ---------------------------
func (u *User) GetFollowers() []*User {
	client := getTwitterClient()
	surveyUserfollowers, _, err := client.Followers.List(&twitter.FollowerListParams{ScreenName: u.ScreenName})
	if err != nil {
		fmt.Println("[Error]: " + err.Error())
		os.Exit(1)
	}

	var followers []*User
	for _, v := range surveyUserfollowers.Users {
		followers = append(followers, NewUser(v.ScreenName))

		// TODO: Cursorがある場合にページネーションして読み込む
	}
	return followers
}

// ShowInfoAsSurveyUser
// 調査対象のアカウント情報を標準出力する
// ---------------------------
func (u *User) ShowInfoAsSurveyUser() {
	fmt.Println("[Info] 調査対象： " + u.AccountName + "(@" + u.ScreenName + ") フォロワー数" + strconv.Itoa(u.FollowerCount) + "人")
}

// ShowInfoAsFollowerUser
// フォロワーのアカウント情報を標準出力する
// ---------------------------
func (u *User) ShowInfoAsFollowerUser(rank int) {
	fmt.Println("[Info] 人気フォロワー" + strconv.Itoa(rank) + "位： " + u.AccountName + "(@" + u.ScreenName + ") フォロワー数" + strconv.Itoa(u.FollowerCount) + "人")
}

// getTwitterClient
// Twitter Clintの取得
// ---------------------------
// TODO: 別のファイルに移動する
func getTwitterClient() *twitter.Client {
    // TODO: シークレットな値は環境変数にする
	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}
