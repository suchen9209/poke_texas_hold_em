package controllers

import (
	"poke/models"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

var langTypes []string // Languages that are supported.

func init() {
	// Initialize language type list.
	lans, _ := beego.AppConfig.String("lang_types")
	langTypes = strings.Split(lans, "|")

	// Load locale files according to language types.
	for _, lang := range langTypes {

		logs.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			logs.Error("Fail to set message file:", err)
			return
		}
	}
}

// baseController represents base router for all other app routers.
// It implemented some methods for the same implementation;
// thus, it will be embedded into other routers.
type baseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (b *baseController) Prepare() {
	// Reset language option.
	b.Lang = "" // This field is from i18n.Locale.

	// 1. Get language information from 'Accept-Language'.
	al := b.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			b.Lang = al
		}
	}

	// 2. Default language is English.
	if len(b.Lang) == 0 {
		b.Lang = "en-US"
	}

	// Set template level language option.
	b.Data["Lang"] = b.Lang
}

// AppController handles the welcome screen that allows user to pick a technology and username.
type AppController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for AppController.
func (a *AppController) Get() {
	a.TplName = "welcome.html"
}

// Join method handles POST requests for AppController.
func (a *AppController) Join() {
	// Get form value.
	uname := a.GetString("uname")

	// Check valid.
	if len(uname) == 0 {
		a.Redirect("/", 302)
		return
	}

	a.Redirect("/ws?uname="+uname, 302)

}

func (a *AppController) GreedIsGood() {

	uname := a.GetString("uname")
	u, _ := models.CheckUser(uname)
	models.ChangeUserPoint(u.Id, 1000)

	a.TplName = "greedisgood.html"
}
