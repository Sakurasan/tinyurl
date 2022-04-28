package tinyurl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tinyurl/pkg/base62"
	"tinyurl/pkg/config"
	data "tinyurl/pkg/db"
	"tinyurl/pkg/log"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/validator"
	lru "github.com/hashicorp/golang-lru"
	"gorm.io/gorm"
)

type Param struct {
	Code     string `json:"code,omitempty"`
	Msg      string `json:"msg,omitempty"`
	LongUrl  string `json:"longUrl,omitempty" ` //validate:"required
	ShortUrl string `json:"shortUrl,omitempty"`
}

func TinyurlHandler(c flamego.Context, logger *log.Logger, db *gorm.DB) string {
	logger.Info(c.Request().Context(), time.Now().Format("2006-01-02 15:04:05"))
	var originurl string
	if len(c.Request().URL.RawQuery) > 1 {
		originurl = c.Param("url") + "?" + c.Request().URL.RawQuery
	} else {
		originurl = c.Param("url")
	}
	return fmt.Sprintf(
		"TinyUrl , %s to %s",
		originurl,
		base62.TinyUrl(originurl),
	)
}

func TinyUrl(c flamego.Context, w http.ResponseWriter, form Param, errs binding.Errors, cfg *config.Config, logger *log.Logger, db *gorm.DB, l *lru.Cache) {
	if len(errs) > 0 {
		var err error
		switch errs[0].Category {
		case binding.ErrorCategoryValidation:
			err = errs[0].Err.(validator.ValidationErrors)[0]
		default:
			err = errs[0].Err
		}
		logger.Error(c.Request().Context(), err.Error())
		c.ResponseWriter().WriteHeader(http.StatusBadRequest)
		_, _ = c.ResponseWriter().Write([]byte(fmt.Sprintf("Oops! Error occurred: %v", err)))
		return
	}
	if !strings.HasPrefix(form.LongUrl, "http://") || !strings.HasPrefix(form.LongUrl, "https://") {
		form.LongUrl = "http://" + form.LongUrl
	}
	logger.Info(c.Request().Context(), form.LongUrl)
	var tm = data.TinyUrl{}
	tm.LongUrl = form.LongUrl
	tm.ShortUrl = base62.TinyUrl(form.LongUrl + "?" + time.Now().String())
	if err := db.Create(&tm).Error; err != nil {
		logger.Info(c.Request().Context(), err.Error())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": http.StatusBadGateway,
			"msg":  "please try again",
		})
	}
	l.Add(tm.ShortUrl, tm.LongUrl)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":     http.StatusOK,
		"msg":      "success",
		"shortUrl": cfg.Server.Domain + tm.ShortUrl,
	})

}

func TinyUrlTo(c flamego.Context, logger *log.Logger, db *gorm.DB, l *lru.Cache, cfg *config.Config) {
	tiny := c.Param("url")
	if v, ok := l.Get(tiny); ok {
		c.Redirect(v.(string), http.StatusFound)
		return
	}
	var tm data.TinyUrl
	db.Where("short_url = ?", tiny).First(&tm)
	l.Add(tiny, tm.LongUrl)
	c.Redirect(cfg.Server.Domain+tm.LongUrl, http.StatusFound)
}

// type tinyurl interface {
// 	CreatreTinyurl()
// 	SetTinyurl()
// }

// type tiny struct {
// 	db *gorm.DB
// }

// func newTiny() *tiny {
// 	return new(tiny)
// }

// func (t *tiny) CreatreTinyurl() {

// 	t.db.Create()
// }
