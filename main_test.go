package main

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/gomodule/oauth1/oauth"
)

type MockClient struct {
	client oauth.Client
}

func (mc *MockClient) ReqGet(credentials *oauth.Credentials, urlStr string) ([]byte, error) {
	//Return test data
	var responseBody []byte

	switch urlStr {
	case "https://api.twitter.com/1.1/account/verify_credentials.json":
		//TODO response body should be recorded beforehand, and save in a test data file
		responseBody = []byte(`{"id":1158036863057158144,"id_str":"1158036863057158144","name":"Miki.masumomo","screen_name":"m_miki0108","location":"Tokyo-to, Japan","description":"\u7b4b\u30c8\u30ec\u3068\u304a\u7d75\u304b\u304d\u304c\u8da3\u5473\u306e\u30d5\u30ea\u30fc\u30e9\u30f3\u30b9\u30a8\u30f3\u30b8\u30cb\u30a2\u3067\u3059|CC11 in @codechrysalis |\u307e\u308b\u3067\u8a71\u305b\u306cTOEIC900|\u884c\u672b\u30bb\u30df\u30b3\u30ed\u30f3\u304c\u3042\u308b\u65b9\u304c\u597d\u307f|\u7d4c\u9a13\u9577\u3044\u306e\u306fJava\u3068JavaScript|\u57fa\u672c\u5fdc\u7528\u9ad8\u5ea6\uff08DB,NW,SC\uff09\u4fdd\u6301 | Golang\u52c9\u5f37\u4e2d","url":"https:\/\/t.co\/xenkzq0x7f","entities":{"url":{"urls":[{"url":"https:\/\/t.co\/xenkzq0x7f","expanded_url":"https:\/\/github.com\/masumomo","display_url":"github.com\/masumomo","indices":[0,23]}]},"description":{"urls":[]}},"protected":false,"followers_count":142,"friends_count":115,"listed_count":3,"created_at":"Sun Aug 04 15:27:58 +0000 2019","favourites_count":1727,"utc_offset":null,"time_zone":null,"geo_enabled":false,"verified":false,"statuses_count":190,"lang":null,"status":{"created_at":"Thu Aug 13 12:23:33 +0000 2020","id":1293885904738586624,"id_str":"1293885904738586624","text":"\u6b63\u76f4RUN\u3068CMD\u3068ENTRYPOINT\u4f7f\u3044\u5206\u3051\u5206\u304b\u3063\u3066\u306a\u304b\u3063\u305f\u3051\u3069\u3001go\u8a00\u8a9e\u3092docker\u3067\u52d5\u304b\u305d\u3046\u3063\u3066\u6642\u306b\u898b\u3064\u3051\u305f\u8a18\u4e8b\u304c\u308f\u304b\u308a\u3084\u3059\u304b\u3063\u305f\u30e1\u30e2\nhttps:\/\/t.co\/po4cF1VxaJ\nRUN\uff1a\u30a4\u30f3\u30b9\u30c8\u30fc\u30eb\u7cfb\u30b3\u30de\u30f3\u30c9\nCM\u2026 https:\/\/t.co\/ogSIv4s1sf","truncated":true,"entities":{"hashtags":[],"symbols":[],"user_mentions":[],"urls":[{"url":"https:\/\/t.co\/po4cF1VxaJ","expanded_url":"https:\/\/goinbigdata.com\/docker-run-vs-cmd-vs-entrypoint\/","display_url":"goinbigdata.com\/docker-run-vs-\u2026","indices":[73,96]},{"url":"https:\/\/t.co\/ogSIv4s1sf","expanded_url":"https:\/\/twitter.com\/i\/web\/status\/1293885904738586624","display_url":"twitter.com\/i\/web\/status\/1\u2026","indices":[117,140]}]},"source":"\u003ca href=\"https:\/\/mobile.twitter.com\" rel=\"nofollow\"\u003eTwitter Web App\u003c\/a\u003e","in_reply_to_status_id":null,"in_reply_to_status_id_str":null,"in_reply_to_user_id":null,"in_reply_to_user_id_str":null,"in_reply_to_screen_name":null,"geo":null,"coordinates":null,"place":null,"contributors":null,"is_quote_status":false,"retweet_count":0,"favorite_count":10,"favorited":false,"retweeted":false,"possibly_sensitive":false,"lang":"ja"},"contributors_enabled":false,"is_translator":false,"is_translation_enabled":false,"profile_background_color":"F5F8FA","profile_background_image_url":null,"profile_background_image_url_https":null,"profile_background_tile":false,"profile_image_url":"http:\/\/pbs.twimg.com\/profile_images\/1262374234967257088\/Zh2DBoDs_normal.jpg","profile_image_url_https":"https:\/\/pbs.twimg.com\/profile_images\/1262374234967257088\/Zh2DBoDs_normal.jpg","profile_banner_url":"https:\/\/pbs.twimg.com\/profile_banners\/1158036863057158144\/1578059311","profile_link_color":"1DA1F2","profile_sidebar_border_color":"C0DEED","profile_sidebar_fill_color":"DDEEF6","profile_text_color":"333333","profile_use_background_image":true,"has_extended_profile":false,"default_profile":true,"default_profile_image":false,"following":false,"follow_request_sent":false,"notifications":false,"translator_type":"none","suspended":false,"needs_phone_verification":false}`)
	case "https://api.twitter.com/1.1/statuses/home_timeline.json":
		return nil, errors.New("Not implimented")
	}
	return responseBody, nil
}

func (mc *MockClient) ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	//Return test data
	fmt.Println("Call mock Post function! but it's not implimented")
	return nil, errors.New("Not implimented")
}

var twitterMockClient MockClient = MockClient{
	client: oauth.Client{}}

func TestGetUser(t *testing.T) {
	var actual User
	expected := User{
		Id:         1158036863057158144,
		Name:       "Miki.masumomo",
		ScreenName: "m_miki0108",
	}

	GetUser(&twitterMockClient, nil, &actual) // Don't need credential

	t.Log(actual)
	if actual != expected {
		t.Fatalf("User must be %v", expected)
	}
}
