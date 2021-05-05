package jpx_business_day

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_businessDay_LastUpdateDate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		businessDay *businessDay
		want        time.Time
	}{
		{name: "ゼロ値ならゼロ値を返す", businessDay: &businessDay{}, want: time.Time{}},
		{name: "timeがあればtimeを返す",
			businessDay: &businessDay{lastUpdateDate: time.Date(2021, 5, 5, 6, 29, 0, 0, time.Local)},
			want:        time.Date(2021, 5, 5, 6, 29, 0, 0, time.Local)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.businessDay.LastUpdateDate()
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_businessDay_LastHoliday(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		businessDay *businessDay
		want        time.Time
	}{
		{name: "ゼロ値ならゼロ値を返す", businessDay: &businessDay{}, want: time.Time{}},
		{name: "値があれば値を返す",
			businessDay: &businessDay{lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			want:        time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.businessDay.LastHoliday()
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_businessDay_IsHoliday(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		businessDay *businessDay
		arg         time.Time
		want        bool
	}{
		{name: "土曜日はtrue",
			businessDay: &businessDay{lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			want:        true},
		{name: "lastHoliday以降でも土曜日はtrue",
			businessDay: &businessDay{lastHoliday: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			want:        true},
		{name: "日曜日はtrue",
			businessDay: &businessDay{lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 2, 0, 0, 0, 0, time.Local),
			want:        true},
		{name: "lastHoliday以降でも日曜日はtrue",
			businessDay: &businessDay{lastHoliday: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 2, 0, 0, 0, 0, time.Local),
			want:        true},
		{name: "lastHoliday以降の平日はfalse",
			businessDay: &businessDay{lastHoliday: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local),
			want:        false},
		{name: "lastHolidayがholidaysにあればtrue",
			businessDay: &businessDay{
				holidays: map[time.Time]string{
					time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
				},
				lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local),
			want: true},
		{name: "lastHoliday以前の平日がholidaysにあればtrue",
			businessDay: &businessDay{
				holidays: map[time.Time]string{
					time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local):   "こどもの日",
					time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
				},
				lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:  time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local),
			want: true},
		{name: "lastHoliday以前の平日がholidaysになければfalse",
			businessDay: &businessDay{
				holidays: map[time.Time]string{
					time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local):   "こどもの日",
					time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
				},
				lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:  time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
			want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.businessDay.IsHoliday(test.arg)
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_businessDay_IsBusinessDay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		businessDay *businessDay
		arg         time.Time
		want        bool
	}{
		{name: "土曜日はfalse",
			businessDay: &businessDay{lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			want:        false},
		{name: "lastHoliday以降でも土曜日はfalse",
			businessDay: &businessDay{lastHoliday: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			want:        false},
		{name: "日曜日はfalse",
			businessDay: &businessDay{lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 2, 0, 0, 0, 0, time.Local),
			want:        false},
		{name: "lastHoliday以降でも日曜日はfalse",
			businessDay: &businessDay{lastHoliday: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 2, 0, 0, 0, 0, time.Local),
			want:        false},
		{name: "lastHoliday以降の平日はtrue",
			businessDay: &businessDay{lastHoliday: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:         time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local),
			want:        true},
		{name: "lastHolidayがholidaysにあればfalse",
			businessDay: &businessDay{
				holidays: map[time.Time]string{
					time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
				},
				lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local),
			want: false},
		{name: "lastHoliday以前の平日がholidaysにあればfalse",
			businessDay: &businessDay{
				holidays: map[time.Time]string{
					time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local):   "こどもの日",
					time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
				},
				lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:  time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local),
			want: false},
		{name: "lastHoliday以前の平日がholidaysになければtrue",
			businessDay: &businessDay{
				holidays: map[time.Time]string{
					time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local):   "こどもの日",
					time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
				},
				lastHoliday: time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local)},
			arg:  time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
			want: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.businessDay.IsBusinessDay(test.arg)
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_businessDay_Refresh_OK(t *testing.T) {
	t.Parallel()
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(source))
	}))
	bd := &businessDay{url: serv.URL}
	err := bd.Refresh(context.Background())
	if err != nil {
		t.Errorf("%s error: %+v\n", t.Name(), err)
	}

	wantHoliday := map[time.Time]string{
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local):   "元日",
		time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local):   "休業日",
		time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local):   "休業日",
		time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local):  "成人の日",
		time.Date(2021, 2, 11, 0, 0, 0, 0, time.Local):  "建国記念の日",
		time.Date(2021, 2, 23, 0, 0, 0, 0, time.Local):  "天皇誕生日",
		time.Date(2021, 3, 20, 0, 0, 0, 0, time.Local):  "春分の日",
		time.Date(2021, 4, 29, 0, 0, 0, 0, time.Local):  "昭和の日",
		time.Date(2021, 5, 3, 0, 0, 0, 0, time.Local):   "憲法記念日",
		time.Date(2021, 5, 4, 0, 0, 0, 0, time.Local):   "みどりの日",
		time.Date(2021, 5, 5, 0, 0, 0, 0, time.Local):   "こどもの日",
		time.Date(2021, 7, 22, 0, 0, 0, 0, time.Local):  "海の日",
		time.Date(2021, 7, 23, 0, 0, 0, 0, time.Local):  "スポーツの日",
		time.Date(2021, 8, 8, 0, 0, 0, 0, time.Local):   "山の日",
		time.Date(2021, 8, 9, 0, 0, 0, 0, time.Local):   "振替休日",
		time.Date(2021, 9, 20, 0, 0, 0, 0, time.Local):  "敬老の日",
		time.Date(2021, 9, 23, 0, 0, 0, 0, time.Local):  "秋分の日",
		time.Date(2021, 11, 3, 0, 0, 0, 0, time.Local):  "文化の日",
		time.Date(2021, 11, 23, 0, 0, 0, 0, time.Local): "勤労感謝の日",
		time.Date(2021, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local):   "元日",
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local):   "休業日",
		time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local):   "休業日",
		time.Date(2022, 1, 10, 0, 0, 0, 0, time.Local):  "成人の日",
		time.Date(2022, 2, 11, 0, 0, 0, 0, time.Local):  "建国記念の日",
		time.Date(2022, 2, 23, 0, 0, 0, 0, time.Local):  "天皇誕生日",
		time.Date(2022, 3, 21, 0, 0, 0, 0, time.Local):  "春分の日",
		time.Date(2022, 4, 29, 0, 0, 0, 0, time.Local):  "昭和の日",
		time.Date(2022, 5, 3, 0, 0, 0, 0, time.Local):   "憲法記念日",
		time.Date(2022, 5, 4, 0, 0, 0, 0, time.Local):   "みどりの日",
		time.Date(2022, 5, 5, 0, 0, 0, 0, time.Local):   "こどもの日",
		time.Date(2022, 7, 18, 0, 0, 0, 0, time.Local):  "海の日",
		time.Date(2022, 8, 11, 0, 0, 0, 0, time.Local):  "山の日",
		time.Date(2022, 9, 19, 0, 0, 0, 0, time.Local):  "敬老の日",
		time.Date(2022, 9, 23, 0, 0, 0, 0, time.Local):  "秋分の日",
		time.Date(2022, 10, 10, 0, 0, 0, 0, time.Local): "スポーツの日",
		time.Date(2022, 11, 3, 0, 0, 0, 0, time.Local):  "文化の日",
		time.Date(2022, 11, 23, 0, 0, 0, 0, time.Local): "勤労感謝の日",
		time.Date(2022, 12, 31, 0, 0, 0, 0, time.Local): "休業日",
	}
	wantLastHoliday := time.Date(2022, 12, 31, 0, 0, 0, 0, time.Local)
	wantLastUpdateDate := time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local)

	if !reflect.DeepEqual(wantHoliday, bd.holidays) || !reflect.DeepEqual(wantLastHoliday, bd.lastHoliday) || !reflect.DeepEqual(wantLastUpdateDate, bd.lastUpdateDate) {
		t.Errorf("%s error\nwant: %+v, %+v, %+v\ngot: %+v, %+v, %+v\n", t.Name(), wantHoliday, wantLastHoliday, wantLastUpdateDate, bd.holidays, bd.lastHoliday, bd.lastUpdateDate)
	}
}

func Test_businessDay_Refresh_Not_OK(t *testing.T) {
	t.Parallel()
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(source))
	}))
	bd := &businessDay{url: serv.URL}
	err := bd.Refresh(context.Background())
	if !errors.Is(err, NotOKStatusError) {
		t.Errorf("%s error: %+v\n", t.Name(), err)
	}
}

var source = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" lang="ja" xml:lang="ja">
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width,initial-scale=1.0" />
<meta http-equiv="Content-Style-Type"  content="text/css" />
<meta http-equiv="Content-Script-Type" content="text/javascript" />
<meta http-equiv="Pragma" content="no-cache">
<meta http-equiv="Cache-Control" content="no-cache">
<meta name="copyright" content="(C) Japan Exchange Group, Inc." />
<meta name="description" content="日本取引所グループ（JPX）は、東京証券取引所、大阪取引所、東京商品取引所等を運営する取引所グループです。総合的なサービス提供を行うことで、市場利用者の方々にとって、より安全で利便性の高い取引の場を提供します。" />
<meta name="keywords" content="日本取引所グループ,JPX,東京証券取引所グループ,大阪取引所,東証,大証" />
<meta property="og:title" content="営業時間・休業日一覧 | 日本取引所グループ">
<meta property="og:type" content="article">
<meta property="og:description" content="日本取引所グループは、東京証券取引所、大阪取引所、東京商品取引所等を運営する取引所グループです。">
<meta property="og:image" content="https://www.jpx.co.jp/common/images/other/nlsgeu000000pud7-img/ogp.jpg">
<meta property="og:site_name" content="日本取引所グループ">
<meta property="fb:admins" content="175272119257459">
<link href="/common/images/icon/nlsgeu000000oie0-img/favicon.ico" rel="shortcut icon" type="image/x-icon" />
<title>営業時間・休業日一覧 | 日本取引所グループ</title>
      <link href="/common/stylesheets/reset.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/layout.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/parts.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/niceforms-default.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/jquery.jscrollpane.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/prettyPhoto-parts.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/list_add.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/print.css" rel="stylesheet" type="text/css" media="print" />

      <link href="/common/stylesheets/probo.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/jquery.bxslider.min.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/style_add.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/style_corporate.css" rel="stylesheet" type="text/css" media="all" />

      <link href="/common/stylesheets/lity.css" rel="stylesheet" type="text/css" media="all" />

        <script type="text/javascript" src="/public/javascripts/jquery.js"></script>
<script type="text/javascript" src="/common/javascripts/niceforms.js"></script>
<script type="text/javascript" src="/common/javascripts/nlsgeu000003yqav-att/lity.js"></script>

<script type="text/javascript" src="/common/javascripts/tvdivq0000001ub2-att/link.js"></script>

<script type="text/javascript" src="/common/javascripts/tvdivq000000dlac-att/jquery.prettyPhoto.js"></script>
<script type="text/javascript" src="/common/javascripts/jquery.mousewheel.js"></script>
<script type="text/javascript" src="/common/javascripts/tvdivq0000004tf9-att/jquery.jscrollpane.min.js"></script>

<script type="text/javascript" src="/common/javascripts/tvdivq0000016ade-att/jquery.easing.1.3.js"></script>

<script type="text/javascript" src="/common/javascripts/tvdivq0000016abs-att/jquery.aslider.min.js"></script>

<script type="text/javascript" src="/common/javascripts/nlsgeu0000038o5u-att/jquery.bxslider.min.js"></script>

<script type="text/javascript" src="/common/javascripts/tvdivq0000025k8e-att/iframe.js"></script>

<script type="text/javascript" src="/common/javascripts/tvdivq00000032wq-att/heightLine.js"></script>

<script type="text/javascript" src="/common/javascripts/nlsgeu000003kowk-att/fixed_midashi.js"></script>

<script type="text/javascript" src="/common/javascripts/nlsgeu000001k7wj-att/fix-table-header.js"></script>

<script type="text/javascript" src="/common/javascripts/nlsgeu0000038o8u-att/add-action.js"></script>

<script>
var UA = (function() {
	var ua = window.navigator.userAgent;
	var reIeU7 = new RegExp('msie [1-7]\\.');
	var ieU7flg = ua.toLowerCase().search(reIeU7) > 0 ? true : false;
	var iphoneflg = ua.toLowerCase().indexOf('iphone') > 0 ? true : false ;
	var ipadflg = ua.toLowerCase().indexOf('ipad') > 0 ? true : false ;
	var androidflg = ua.toLowerCase().indexOf('android') > 0 ? true : false ;
	var mobileflg = ua.toLowerCase().indexOf('mobile') > 0 ? true : false

	return {
		isIeU7 : function() {
			return ieU7flg;
		}
		, isSp : function() {
			return iphoneflg || ipadflg || androidflg ? true : false;
		}
	};
})();
</script>

<script src="//cdn1.readspeaker.com/script/6483/webReader/webReader.js?pids=wr&amp;forceAdapter=ioshtml5&amp;disable=translation,lookup" type="text/javascript"></script>
<!-- Facebook Pixel Code -->
<script>
  !function(f,b,e,v,n,t,s)
  {if(f.fbq)return;n=f.fbq=function(){n.callMethod?
  n.callMethod.apply(n,arguments):n.queue.push(arguments)};
  if(!f._fbq)f._fbq=n;n.push=n;n.loaded=!0;n.version='2.0';
  n.queue=[];t=b.createElement(e);t.async=!0;
  t.src=v;s=b.getElementsByTagName(e)[0];
  s.parentNode.insertBefore(t,s)}(window, document,'script',
  'https://connect.facebook.net/en_US/fbevents.js');
  fbq('init', '191019531472708');
  fbq('track', 'PageView');
</script>
<noscript><img height="1" width="1" style="display:none"
  src="https://www.facebook.com/tr?id=191019531472708&ev=PageView&noscript=1"
/></noscript>
<!-- End Facebook Pixel Code -->

</head>

<body class="corporate">

<!-- body_prepend -->
<!-- Google Tag Manager -->
<link rel="stylesheet" type="text/css" href="/common/stylesheets/cookieconsent_customize.css" />
<script src="//cdnjs.cloudflare.com/ajax/libs/cookieconsent2/3.1.0/cookieconsent.min.js"></script>
<script src="/common/javascripts/nlsgeu000003vo9u.js"></script>
<noscript><iframe src="//www.googletagmanager.com/ns.html?id=GTM-MS2WRF"
height="0" width="0" style="display:none;visibility:hidden"></iframe></noscript>
<script>
function startGtm() {
(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
'//www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','GTM-MS2WRF');
}
</script>
<!-- End Google Tag Manager -->
<div id="wrapper-area">
<div id="header-area">
<div id="header-areaIn">
<div class="header-menu">
<ul>
<li><a href="/">JPX トップページへ</a></li>
<li><a href="/corporate/about-jpx/access/index.html">アクセス</a></li>
<li><a href="/contact/index.html">お問合せ</a></li>
</ul>
</div>
<div class="header-btn">
<ul class="language-btn">
<li><a href="/english/corporate/">English</a></li>
<li><a href="/chinese/corporate/jpx-profile/">中文</a></li>
</ul>
<ul class="fontsize-btn">
<li>文字サイズ</li>
<li><a href="javaScript:void(0)" data-size="small">小</a></li>
<li class="act"><a href="javaScript:void(0)" data-size="medium">中</a></li>
<li><a href="javaScript:void(0)" data-size="large">大</a></li>
</ul>
<div class="search-btn">
  <form action="/search.html">
    <input type="text" id="q2" class="search-input" name="q" value="" placeholder="検索キーワード">
    <div class="search-input-btn"><input type="image" src="/common/images/other/tvdivq000000klrc-img/btn-search.png" name="" /></div>
  </form>
</div>
</div>
</div>
</div>


	<div class="bread-crumb-box">
		<div class="bread-crumb">
			<ol>

						<li><a href="/corporate/index.html">JPXについて</a></li>

						<li><a href="/corporate/about-jpx/index.html">会社情報</a></li>

						<li><a href="/corporate/about-jpx/calendar/index.html">営業時間・休業日一覧</a></li>

			</ol>
		</div>
	</div>


  <div id="menu-area">
 <div id="menu-areaIn" style="height: 927px;">
<p id="site-logo">
  <a href="/corporate">
    <img class="onlypc" src="/common/images/other/logo_corporate.png" width="207" alt="JPX 日本取引所グループ" title="JPX 日本取引所グループ" />
    <img class="onlysp" src="/common/images/other/logo_corporate_sp.png" alt="JPX 日本取引所グループ" title="JPX 日本取引所グループ">
  </a>
  <a class="onlysp" id="spmenu-open">
    <img class="onlysp" src="/common/images/other/menu_corporate_sp.png" alt="JPX 日本取引所グループ" title="JPX 日本取引所グループ">
  </a>
</p>

<div id="main-menu">
  <div class="search-btn onlysp">
    <form method="get" name="" id="" action="/search.html">
      <input type="text" id="q2" class="search-input" name="q" value="" placeholder="検索キーワード" />
      <div class="search-input-btn">
        <input type="image" src="/common/images/other/tvdivq000000klrc-img/btn-search.png" name="" value="" />
      </div>
    </form>
  </div>
  <p id="menu-btn-slide" class="onlypc menu-btn-open"><a href="javaScript:void(0)">MENU</a></p>
  <p id="sp-menu-btn-slide" class="menu-btn-open onlysp"><a>MENU</a></p>
  <div id="main-menuIn">
    <h2><a href="/corporate/ceo-message/index.html" data-target="menutvdivq0000006o89">グループCEOごあいさつ</a></h2>
<h2><a href="/corporate/news/news-releases/index.html" data-target="menutvdivq000000zn8c">JPXからのお知らせ</a></h2>
<h2><a href="/corporate/about-jpx/index.html" data-target="menun3i7740000001jya">会社情報</a></h2>
<h2><a href="/corporate/governance/index.html" data-target="menun3i7740000001nfr">ガバナンス／リスク管理</a></h2>
<h2><a href="/corporate/investor-relations/index.html" data-target="menutvdivq000000lbh5">株主・投資家情報（IR）</a></h2>
<h2><a href="/corporate/sustainability/index.html" data-target="menunlsgeu0000036fhk">サステナビリティ</a></h2>
<h2><a href="/corporate/research-study/index.html" data-target="menutvdivq000000af0z">調査・研究／政策提言</a></h2>
<h2><a href="/corporate/events-pr/index.html" data-target="menun3i7740000001jte">イベント・PR</a></h2>
<div class="main-menu-sub">
<h2><a href="/">JPX トップページへ</a></h2>
</div>
<div class="main-menu-sub">
<ul class="main-menu-icon">
<li><a href="https://twitter.com/JPX_official" rel="external"><img src="/common/images/icon/icon_menu_twitter.png" alt="Twitter"></a></li>
<li><a href="https://www.facebook.com/JapanExchangeGroup" rel="external"><img src="/common/images/icon/icon_menu_facebook.png" alt="Facebook"></a></li>
<li><a href="https://www.youtube.com/channel/UCnZA74T8a8dEbavWRq8F2nA" rel="external"><img src="/common/images/icon/icon_menu_youtube.png" alt="Youtube"></a></li>
<li><a href="https://www.instagram.com/jpx_official/" rel="external"><img src="/common/images/icon/icon_menu_instagram.png" alt="Instagram"></a></li>
</ul>
<ul class="main-menu-link">
<li><a href="/learning/social-media/index.html">ソーシャルメディア一覧</a></li>
<li><a href="/learning/mail-magazine/index.html">メールマガジン</a></li>
</ul>
</div>
</div>
<div id="sp-main-menuIn">
<h2><a href="/corporate/ceo-message/index.html" data-target="menutvdivq0000006o89">グループCEOごあいさつ</a></h2>
<h2><a href="/corporate/news/news-releases/index.html" data-target="menutvdivq000000zn8c">JPXからのお知らせ</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/news/news-releases/index.html">ニュース一覧</a></li>
<li><a href="/corporate/news/monthly-headline/index.html">JPXマンスリー・ヘッドライン</a></li>
<li><a href="/corporate/news/press-conference/index.html">CEO定例記者会見</a></li>
</ul>
</div>
<h2><a href="/corporate/about-jpx/index.html" data-target="menun3i7740000001jya">会社情報</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/about-jpx/business/index.html">事業紹介</a></li>
<li><a href="/corporate/about-jpx/philosophy/index.html">企業理念</a></li>
<li><a href="/corporate/about-jpx/jpx-logo/index.html">コーポレートロゴ</a></li>
<li><a href="/corporate/about-jpx/profile/index.html">会社概要</a></li>
<li><a href="/corporate/about-jpx/officer/index.html">役員一覧</a></li>
<li><a href="/corporate/about-jpx/organization/index.html">組織図</a></li>
<li><a href="/corporate/about-jpx/history/index.html">沿革</a></li>
<li><a href="/corporate/about-jpx/calendar/index.html">営業時間・休業日一覧</a></li>
<li><a href="/corporate/about-jpx/access/index.html">アクセス</a></li>
<li><a href="/corporate/about-jpx/recruit/index.html">採用情報</a></li>
</ul>
</div>
<h2><a href="/corporate/governance/index.html" data-target="menun3i7740000001nfr">ガバナンス／リスク管理</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/governance/charter/index.html">企業行動憲章</a></li>
<li><a href="/corporate/governance/policy/index.html">コーポレート・ガバナンス</a></li>
<li><a href="/corporate/governance/self-regulation/index.html">自主規制業務の適正な体制整備</a></li>
<li><a href="/corporate/governance/compliance/index.html">コンプライアンス・プログラム</a></li>
<li><a href="/corporate/governance/internal-control/index.html">内部統制システム構築の基本方針</a></li>
<li><a href="/corporate/governance/risk/index.html">リスク管理</a></li>
<li><a href="/corporate/governance/security/index.html">情報セキュリティ</a></li>
<li><a href="/corporate/governance/principle/index.html">「不祥事予防のプリンシプル」の対応状況</a></li>
</ul>
</div>
<h2><a href="/corporate/investor-relations/index.html" data-target="menutvdivq000000lbh5">株主・投資家情報（IR）</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/investor-relations/management/index.html">経営情報</a></li>
<li><a href="/corporate/investor-relations/individual/index.html">個人投資家の皆様へ</a></li>
<li><a href="/corporate/investor-relations/ir-library/index.html">IR資料室</a></li>
<li><a href="/corporate/investor-relations/financials/index.html">業績・財務</a></li>
<li><a href="/corporate/investor-relations/shareholders/index.html">株主・株式情報</a></li>
<li><a href="/corporate/investor-relations/ir-calendar/index.html">IRカレンダー</a></li>
<li><a href="/corporate/investor-relations/ir-mail/index.html">IRメール配信サービス</a></li>
<li><a href="/corporate/investor-relations/ir-faq/index.html">IRに関するよくあるご質問（FAQ)</a></li>
</ul>
</div>
<h2><a href="/corporate/sustainability/index.html" data-target="menunlsgeu0000036fhk">サステナビリティ</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/sustainability/our-sustainability/index.html">JPXの考えるサステナビリティ</a></li>
<li><a href="/corporate/sustainability/esg-investment/index.html">ESG投資の普及に向けた取組み</a></li>
<li><a href="/corporate/sustainability/jpx-esg/index.html">JPXのESG情報</a></li>
<li><a href="/corporate/sustainability/esgknowledgehub/index.html">JPX ESG Knowledge Hub</a></li>
<li><a href="/corporate/sustainability/news-events/index.html">関連ニュース・イベント</a></li>
</ul>
</div>
<h2><a href="/corporate/research-study/index.html" data-target="menutvdivq000000af0z">調査・研究／政策提言</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/research-study/working-paper/index.html">JPXワーキング・ペーパー</a></li>
<li><a href="/corporate/research-study/research-group/index.html">日本取引所グループ金融商品取引法研究会</a></li>
<li><a href="/corporate/research-study/suggestions/index.html">JPX金融資本市場ワークショップからの提言</a></li>
<li><a href="/corporate/research-study/derivatives/index.html">デリバティブ投資家層の裾野拡大に向けた勉強会</a></li>
<li><a href="/corporate/research-study/dlt/index.html">業界連携型DLT実証実験</a></li>
<li><a href="/corporate/research-study/research-archives/macro-group/index.html">過去の各種研究会</a></li>
<li><a href="/corporate/research-study/system-failure/index.html">システム障害に係る「再発防止策検討協議会」</a></li>
</ul>
</div>
<h2><a href="/corporate/events-pr/index.html" data-target="menun3i7740000001jte">イベント・PR</a><span></span>
</h2>
<div class="sp-sub-menu">
<ul>
<li><a href="/corporate/events-pr/ceremony/index.html">大納会・大発会</a></li>
<li><a href="/corporate/events-pr/concert/index.html">JPXコンサート</a></li>
<li><a href="/corporate/events-pr/140years/index.html">株式取引所開設140周年</a></li>
</ul>
</div>
<div class="main-menu-sub">
<h3><a href="/">JPX トップページへ</a></h3>
</div>
<div class="main-menu-sub">
<ul class="main-menu-icon">
<li><a href="https://twitter.com/JPX_official" rel="external"><img src="/common/images/icon/icon_menu_twitter.png" alt="Twitter"></a></li>
<li><a href="https://www.facebook.com/JapanExchangeGroup" rel="external"><img src="/common/images/icon/icon_menu_facebook.png" alt="Facebook"></a></li>
<li><a href="https://www.youtube.com/channel/UCnZA74T8a8dEbavWRq8F2nA" rel="external"><img src="/common/images/icon/icon_menu_youtube.png" alt="Youtube"></a></li>
<li><a href="https://www.instagram.com/jpx_official/" rel="external"><img src="/common/images/icon/icon_menu_instagram.png" alt="Instagram"></a></li>
</ul>
<ul class="main-menu-link">
<li><a href="/learning/social-media/index.html">ソーシャルメディア一覧</a></li>
<li><a href="/learning/mail-magazine/index.html">メールマガジン</a></li>
</ul>
</div>
  </div>

  <div id="sub-menu">

									<h2 class="menu-title"><a href="/corporate/about-jpx/index.html">会社情報</a></h2>

										<h3 class="sub-title"><a href="/corporate/about-jpx/business/index.html"  data-target="sub-menutvdivq000000v75t">事業紹介</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/philosophy/index.html"  data-target="sub-menutvdivq0000006p0x">企業理念</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/jpx-logo/index.html"  data-target="sub-menutvdivq0000007ijd">コーポレートロゴ</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/profile/index.html"  data-target="sub-menutvdivq0000007ru1">会社概要</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/officer/index.html"  data-target="sub-menutvdivq0000007s7n">役員一覧</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/organization/index.html"  data-target="sub-menutvdivq0000007tng">組織図</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/history/index.html"  data-target="sub-menutvdivq0000007u0g">沿革</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/calendar/index.html" class="act-link" data-target="sub-menunlsgeu000002vfaf">営業時間・休業日一覧</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/access/index.html"  data-target="sub-menutvdivq0000008oov">アクセス</a></h3>

										<h3 class="sub-title"><a href="/corporate/about-jpx/recruit/index.html"  data-target="sub-menutvdivq0000008p74">採用情報</a></h3>




  </div>
  <div class="sp-submenu-wrap">
    <ul class="sp-submenus onlysp">
      <li><a href="/corporate/about-jpx/access/index.html">アクセス</a></li>
      <li><a href="/contact/index.html">お問合せ</a></li>
    </ul>
    <ul class="language-btn onlysp">
        <li><a href="/english/corporate/">English</a></li>
        <li><a href="/chinese/corporate/jpx-profile/">中文</a></li>
    </ul>
  </div>
  <p class="menu-close onlysp"><i><img src="/common/images/icon/icon-close-white.png"></i>閉じる</p>
</div>



</div><!-- /menu-areaIn -->

</div><!-- /menu-area -->

  <div id="main-area">
    <div id="main-areaIn">
      <div id="read-area">
  <ul>
    <li>2021/01/07 更新</li>
    <li class="icon-read"><a class="rs_href" rel="nofollow" accesskey="L" href="//app-as.readspeaker.com/cgi-bin/rsent?customerid=6483&amp;lang=ja_jp&amp;readid=readArea&amp;url=" target="_blank" onclick="readpage(this.href, 'xp1'); return false;"><div class="onlypc">このページを音声で聴く</div></a></li>
    <li class="icon-print"><a href="javascript:void(0)" onclick="window.print();return false;">印刷</a></li>
  </ul>
</div>
<div id="xp1" class="rs_preserve rs_skip rs_splitbutton rs_addtools rs_exp"></div>
      <div id="readArea">
        <!-- 大見出し -->
        <div><div class="headline-title-wrap"><h1 class="headline-title"><span>営業時間・休業日一覧</span></h1></div></div>




              <div>
                <div class="tab-submenu-anchor">
                  <ul>

                          <li><a href="#heading_0" class="link-window">営業時間</a></li>

                          <li><a href="#heading_9" class="link-window">休業日一覧</a></li>

                  </ul>
                </div>
              </div>
            <div><h2 class="heading-title-mu" id="heading_0"><span>営業時間</span></h2></div>
          <div><p class="component-text">8時45分～16時45分（月～金、祝日を除く）</p></div>
    <div><h3 class="subhead-title" id="heading_2"><span>東京証券取引所の売買立会時間（現物市場）</span></h3></div>
          <div><p class="component-text">内国株式・外国株式・ETF・ETN・REIT・債券（国債を除く）の立会時間はこちらから。<br/>
<a href="https://www.jpx.co.jp/equities/trading/domestic/01.html" class="link-window">売買立会時（立会時間）</a><br/>
<br/>
国債の立会時間はこちらから。<br/>
<a href="https://www.jpx.co.jp/equities/products/bonds/trading/index.html" class="link-window">国債・売買制度</a></p></div>
    <div><h3 class="subhead-title" id="heading_4"><span>大阪取引所の取引時間（先物・オプション市場）</span></h3></div>
          <div><p class="component-text">先物・オプション取引の立会時間はこちらから。<br/>
<a href="https://www.jpx.co.jp/derivatives/rules/trading-hours/index.html" class="link-window">立会時間</a><br/>
<br/>
ナイト・セッションの対象取引はこちらから。<br/>
<a href="https://www.jpx.co.jp/derivatives/rules/night-session/index.html" class="link-window">ナイト・セッション</a></p></div>

<div>

    <div class="component-annotation">
      <ul>

              <li>先物・オプション市場では、営業日の翌日午前5:30までナイト・セッションを行います。営業日の翌日が休日の場合でも翌日午前5:30まで取引を行います。</li>

      </ul>
    </div>

</div>

    <div><h3 class="subhead-title" id="heading_7"><span>東京商品取引所の立会時間（商品先物市場）</span></h3></div>
          <div><p class="component-text">商品先物（原油・エネルギー）市場の立会時間はこちらから。<br/>
<a href="https://www.tocom.or.jp/jp/market/trading_schedule.html" class="link-blank" rel="external">立会時間と計算区域（東京商品取引所ウェブサイト）</a></p></div>
    <div><h2 class="heading-title" id="heading_9"><span>休業日一覧</span></h2></div>
          <div><p class="component-text">休業日一覧は、国民の祝日に関する法律（祝日法）の改正及びその他祝日に関する特別法の制定等により変更になる場合があります。</p></div>
    <div><h3 class="subhead-title" id="heading_11"><span>2021年</span></h3></div>
          <div><p class="component-text"><div class="component-normal-table">
<table class="overtable">
<tr><th width="50%">日付</th><th width="50%">名称</th></tr>
<tr><td class="a-center">2021/01/01（金）</td><td class="a-center">元日</td></tr>
<tr><td class="a-center">2021/01/02（土）</td><td class="a-center">休業日</td></tr>
<tr><td class="a-center">2021/01/03（日）</td><td class="a-center">休業日</td></tr>
<tr><td class="a-center">2021/01/11（月）</td><td class="a-center">成人の日</td></tr>
<tr><td class="a-center">2021/02/11（木）</td><td class="a-center">建国記念の日</td></tr>
<tr><td class="a-center">2021/02/23（火）</td><td class="a-center">天皇誕生日</td></tr>
<tr><td class="a-center">2021/03/20（土）</td><td class="a-center">春分の日</td></tr>
<tr><td class="a-center">2021/04/29（木）</td><td class="a-center">昭和の日</td></tr>
<tr><td class="a-center">2021/05/03（月）</td><td class="a-center">憲法記念日</td></tr>
<tr><td class="a-center">2021/05/04（火）</td><td class="a-center">みどりの日</td></tr>
<tr><td class="a-center">2021/05/05（水）</td><td class="a-center">こどもの日</td></tr>
<tr><td class="a-center">2021/07/22（木）</td><td class="a-center">海の日</td></tr>
<tr><td class="a-center">2021/07/23（金）</td><td class="a-center">スポーツの日</td></tr>
<tr><td class="a-center">2021/08/08（日）</td><td class="a-center">山の日</td></tr>
<tr><td class="a-center">2021/08/09（月）</td><td class="a-center">振替休日</td></tr>
<tr><td class="a-center">2021/09/20（月）</td><td class="a-center">敬老の日</td></tr>
<tr><td class="a-center">2021/09/23（木）</td><td class="a-center">秋分の日</td></tr>
<tr><td class="a-center">2021/11/03（水）</td><td class="a-center">文化の日</td></tr>
<tr><td class="a-center">2021/11/23（火）</td><td class="a-center">勤労感謝の日</td></tr>
<tr><td class="a-center">2021/12/31（金）</td><td class="a-center">休業日</td></tr>
</table>
</div></p></div>
    <div><h3 class="subhead-title" id="heading_13"><span>2022年</span></h3></div>
          <div><p class="component-text"><div class="component-normal-table">
<table class="overtable">
<tr><th width="50%">日付</th><th width="50%">名称</th></tr>
<tr><td class="a-center">2022/01/01（土）</td><td class="a-center">元日</td></tr>
<tr><td class="a-center">2022/01/02（日）</td><td class="a-center">休業日</td></tr>
<tr><td class="a-center">2022/01/03（月）</td><td class="a-center">休業日</td></tr>
<tr><td class="a-center">2022/01/10（月）</td><td class="a-center">成人の日</td></tr>
<tr><td class="a-center">2022/02/11（金）</td><td class="a-center">建国記念の日</td></tr>
<tr><td class="a-center">2022/02/23（水）</td><td class="a-center">天皇誕生日</td></tr>
<tr><td class="a-center">2022/03/21（月）</td><td class="a-center">春分の日</td></tr>
<tr><td class="a-center">2022/04/29（金）</td><td class="a-center">昭和の日</td></tr>
<tr><td class="a-center">2022/05/03（火）</td><td class="a-center">憲法記念日</td></tr>
<tr><td class="a-center">2022/05/04（水）</td><td class="a-center">みどりの日</td></tr>
<tr><td class="a-center">2022/05/05（木）</td><td class="a-center">こどもの日</td></tr>
<tr><td class="a-center">2022/07/18（月）</td><td class="a-center">海の日</td></tr>
<tr><td class="a-center">2022/08/11（木）</td><td class="a-center">山の日</td></tr>
<tr><td class="a-center">2022/09/19（月）</td><td class="a-center">敬老の日</td></tr>
<tr><td class="a-center">2022/09/23（金）</td><td class="a-center">秋分の日</td></tr>
<tr><td class="a-center">2022/10/10（月）</td><td class="a-center">スポーツの日</td></tr>
<tr><td class="a-center">2022/11/03（木）</td><td class="a-center">文化の日</td></tr>
<tr><td class="a-center">2022/11/23（水）</td><td class="a-center">勤労感謝の日</td></tr>
<tr><td class="a-center">2022/12/31（土）</td><td class="a-center">休業日</td></tr>
</table>
</div></p></div>

    </div><!-- /readArea -->
    </div><!-- /main-areaIn -->
      <div id="footer-area">
  <div class="page-top"><p class="icon-pagetop"><a href="#">ページトップ</a></p></div>
  <div class="footer-areaIn">

      <div class="footer-sitetop"><a href="/">JPX トップページへ</a></div>
      <script>
        $(function(){
          $.ajax({
            url:"/p4pd2n00000024m6.xml",
            type:"GET",
            dataType:"xml",
            async:false,
            timeout:1000,
            error:function() {
              $("#therd_map").html('情報を取得できませんでした。');
            },
            success:function(xml){
              var second_menu_count = 0
              var thard_menu_count = 0;
              var count = 0;
              var html = '';
              html += '<div class="footer-sitemap">';
              html += '<ul class="fs-table">';
              $(xml).find("menu").each(function() {
                if (thard_menu_count != count) {
                  count++
                  return true;
                } else {
                  if (second_menu_count%3 == 0 && second_menu_count != 0) {
                    html += '</ul>';
                    html += '</div>';
                    html += '<div class="footer-sitemap">';
                    html += '<ul class="fs-table">';
                  }
                  second_menu_count++;
                }

                if ($(this).find('second_url').text() != "/news/index.html") {
                  html += '<li class="fs-table-sub">';
                } else {
                  html += '<li class="fs-table-sub fs-list-sub">';
                }
                html += '<h2><a href="'+$(this).find('second_url').text()+'">'+$(this).find('second_name').text()+'</a></h2>';
                html += '<ul>';
                $(this).find("menu_list").find("menu").each(function() {
                  html += '<li><a href="'+$(this).find('third_url').text()+'">'+$(this).find('third_name').text()+'</a></li>';
                  thard_menu_count++;
                });
                html += '</ul>';
                html += '</li>';
              });

              if (second_menu_count%3 != 0) {
                while(second_menu_count%3 != 0) {
                  html += '<li class="fs-table-sub">&nbsp</li>';
                  second_menu_count++;
                }
              }

              html += '</ul>';
              html += '</div>';
              $("#therd_map").append(html);
            }
          });
        });
      </script>
      <div class="footer-sitemap-box" id="therd_map">
      </div>
    <div class="footer-menu">
<ul>
<li class="onlysp"><a href="/corporate/">トップページ</a></li>

<li><a href="/site-updates/index.html" >サイト更新情報</a></li>

<li><a href="/faq/index.html" >よくあるご質問</a></li>

<li><a href="/sitemap/index.html" >サイトマップ</a></li>

<li><a href="/term-of-use/index.html" >サイトのご利用上の注意と免責事項</a></li>

<li><a href="/corporate/governance/security/personal-information/index.html" >個人情報の取扱い</a></li>

<li><a href="/corporate/about-jpx/recruit/index.html" >採用情報</a></li>

<li><a href="/corporate/investor-relations/shareholders/announcement/index.html" >法定公告</a></li>

</ul>
</div>

<div class="footer-copyright">
  <p>&copy; <script type="text/javascript">document.write(new Date().getFullYear());</script> Japan Exchange Group, Inc.</p>
</div>
</div>
</div>

  </div><!-- /main-area -->
</div><!-- /wrapper-area -->

<script type="text/javascript" src="/common/javascripts/add_attribute_gid.js"></script>
<!-- body_append -->
<div id="modal-bg"></div>
</body>
</html>
`
