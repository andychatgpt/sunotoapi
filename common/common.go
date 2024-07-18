package common

//goland:noinspection GoSnakeCaseUsage
import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

//goland:noinspection SpellCheckingInspection
const (
	defaultTimeoutSeconds = 600 // 10 minutes
)

var (
	Client tls_client.HttpClient
)

//goland:noinspection GoUnhandledErrorResult
func init() {

	Client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithTimeoutSeconds(defaultTimeoutSeconds),
		tls_client.WithClientProfile(profiles.Chrome_120),
	}...)

}
