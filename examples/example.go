package main

import (
	"context"
	"fmt"
	"time"

	jbd "gitlab.com/tsuchinaga/jpx-business-day"
)

func main() {
	// 取得
	bd := jbd.NewBusinessDay()

	// 初期化チェック
	if bd.LastUpdateDate().IsZero() {

		// 初期化前なのでリフレッシュしておく
		if err := bd.Refresh(context.Background()); err != nil {
			// エラーハンドリング
			panic(err)
		}

		// Refreshでエラーがなければありえないけど、初期化されていることの確認として
		if bd.LastUpdateDate().IsZero() {
			panic("初期化されていません")
		}
	}

	// 営業日かの確認
	now := time.Now()
	target := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local) // 当年の元旦
	if bd.IsHoliday(target) {
		fmt.Printf("%sはお休みです\n", target.Format("2006/01/02"))
	}
}
