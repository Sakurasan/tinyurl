package route

import (
	"tinyurl/route/tinyurl"

	"github.com/flamego/auth"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func Route(f *flamego.Flame) {
	f.Get("/version", auth.Basic("admin", "admin"), func() string { return "1.1.1" })
	// f.Get("/{url: **, capture: 10}", tinyurl.TinyurlHandler)
	f.Get("/{url: **, capture: 10}", tinyurl.TinyUrlTo)
	f.Post("/api/v1/tiny", binding.JSON(tinyurl.Param{}), tinyurl.TinyUrl)

}

func tinyauth(c flamego.Context) {
	if c.Query("token") == "123" {
		c.Redirect("/signup")
	}
}
