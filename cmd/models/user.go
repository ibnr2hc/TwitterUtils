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
	var followers []*User

	// フォロワーの情報を取得する
	var nextCursor int64 = 0
	var followersList []*User
	for true {
		fmt.Println("[Debug] データ取得中... (Cursor: " + strconv.FormatInt(nextCursor, 10) + ")")
		followersList, nextCursor = u.getFollowers(u.ScreenName, nextCursor, client)
		followers = append(followers, followersList...)

		// フォロワー数が多い場合にカーソルでデータを取得するため再度処理を行う
		if nextCursor != 0 {
			continue
		} else {
			break
		}

	}
	return followers
}

func (u *User) getFollowers(screenName string, cursor int64, client *twitter.Client) ([]*User, int64) {
	var followers []*User
	var followerListParams *twitter.FollowerListParams
	if cursor == 0 {
		followerListParams = &twitter.FollowerListParams{ScreenName: screenName}
	} else {
		followerListParams = &twitter.FollowerListParams{Cursor: cursor}
	}

	surveyUserFollowers, _, err := client.Followers.List(followerListParams)
	if err != nil {
		fmt.Println("[Error]: " + err.Error())
		os.Exit(1)
	}

	for _, v := range surveyUserFollowers.Users {
		followers = append(followers, NewUser(v.ScreenName))
	}

	return followers, surveyUserFollowers.NextCursor
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
