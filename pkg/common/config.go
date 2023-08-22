package common

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"gorm.io/gorm"
)

type Config struct {
	CTX        appctx.AppContext
	DB         *gorm.DB
	TokenMaker core.TokenMakerComponent
	Gin        core.GinComponent
	DBMigrator core.DBMigrationComponent
}

func NewConfig(ctx appctx.AppContext) *Config {
	return &Config{
		CTX:        ctx,
		DB:         ctx.MustGet(PluginDBMain).(core.GormDBComponent).GetDB(),
		TokenMaker: ctx.MustGet(PluginTokenMaker).(core.TokenMakerComponent),
		Gin:        ctx.MustGet(PluginGin).(core.GinComponent),
		DBMigrator: ctx.MustGet(PluginDBMigrator).(core.DBMigrationComponent),
	}
}
