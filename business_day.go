package jpx_business_day

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
)

func NewBusinessDay() (BusinessDay, error) {
	bd := &businessDay{url: "https://www.jpx.co.jp/corporate/about-jpx/calendar/"}
	if err := bd.Refresh(context.Background()); err != nil {
		return nil, err
	}
	return bd, nil
}

type BusinessDay interface {
	IsBusinessDay(target time.Time) bool
	IsHoliday(target time.Time) bool
	Refresh(ctx context.Context) error
	LastHoliday() time.Time
	LastUpdateDate() time.Time
}

type businessDay struct {
	url            string
	holidays       map[time.Time]string
	lastHoliday    time.Time
	lastUpdateDate time.Time
	mtx            sync.Mutex
}

// IsBusinessDay - 営業日かどうか
func (b *businessDay) IsBusinessDay(target time.Time) bool {
	return !b.IsHoliday(target)
}

// IsHoliday - 休日かどうか
func (b *businessDay) IsHoliday(target time.Time) bool {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	// 土曜日、日曜日は常に休み
	if target.Weekday() == time.Saturday || target.Weekday() == time.Sunday {
		return true
	}

	// 取得した最終の休日以降は情報がないので常にfalse
	if target.After(b.lastHoliday) {
		return false
	}

	// 祝日一覧にあれば休日
	targetDate := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.Local)
	_, ok := b.holidays[targetDate]
	return ok
}

var (
	NotOKStatusError = errors.New("not ok status error")
	TimeParseError   = errors.New("time parse error")
)

func (b *businessDay) Refresh(ctx context.Context) (err error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET", b.url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := res.Body.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status is %d: %w", res.StatusCode, NotOKStatusError)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)

	updateDate := regexp.MustCompile(`<li>(\d{4}/\d{2}/\d{2}) 更新</li>`).FindAllStringSubmatch(bodyStr, -1)
	if len(updateDate) < 1 || len(updateDate[0]) < 2 {
		return fmt.Errorf("udpate datetime is not found, %w", TimeParseError)
	}
	update, err := time.ParseInLocation("2006/01/02", updateDate[0][1], time.Local)
	if err != nil {
		return fmt.Errorf("%v, %w", err, TimeParseError)
	}
	b.lastUpdateDate = update

	b.holidays = map[time.Time]string{}
	holidays := regexp.MustCompile(`<tr><td class="a-center">(\d{4}/\d{2}/\d{2})\S+</td><td class="a-center">(\S+)</td></tr>`).FindAllStringSubmatch(bodyStr, -1)
	for _, holiday := range holidays {
		if len(holiday) != 3 {
			continue
		}

		if t, err := time.ParseInLocation("2006/01/02", holiday[1], time.Local); err == nil {
			b.holidays[t] = holiday[2]
			b.lastHoliday = t
		}
	}

	return nil
}

// LastHoliday - 取得した最終の休日
func (b *businessDay) LastHoliday() time.Time {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return b.lastHoliday
}

// LastUpdateDate - 営業日情報を取得しているページの更新日
func (b *businessDay) LastUpdateDate() time.Time {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return b.lastUpdateDate
}
