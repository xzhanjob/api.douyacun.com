package db

import (
	"dyc/internal/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"strings"
	"sync"
)

var db_write *gorm.DB
var DB_ONCE sync.Once

type _db struct {
	gorm.DB
}

const (
	parentSpanGormKey = "opentracingParentSpan"
	spanGormKey       = "opentracingSpan"
)

func NewDB(dsn string) {
	DB_ONCE.Do(func() {
		// videos_t
		var err error
		db_write, err = gorm.Open("mysql", dsn)
		if err != nil {
			logger.Fatal("gorm open %s", err)
		}
		if err = db_write.DB().Ping(); err != nil {
			logger.Fatal(err)
		}
		db_write.SingularTable(true)
		AddGormCallbacks(db_write)
	})
}

func Write(ctx *gin.Context) *gorm.DB {
	return SetSpanToGorm(ctx, db_write)
}

// SetSpanToGorm sets span to gorm settings, returns cloned DB
func SetSpanToGorm(ctx *gin.Context, d *gorm.DB) *gorm.DB {
	if ctx == nil {
		return d
	}
	if parentSpan, ok := ctx.Get("root_span"); ok {
		return d.Set(parentSpanGormKey, parentSpan)
	}
	return d
}

// AddGormCallbacks adds callbacks for tracing, you should call SetSpanToGorm to make them work
func AddGormCallbacks(db *gorm.DB) {
	callbacks := newCallbacks()
	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.Scope)    { c.after(scope, "INSERT") }
func (c *callbacks) beforeQuery(scope *gorm.Scope)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.Scope)     { c.after(scope, "SELECT") }
func (c *callbacks) beforeUpdate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.Scope)    { c.after(scope, "UPDATE") }
func (c *callbacks) beforeDelete(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.Scope)    { c.after(scope, "DELETE") }
func (c *callbacks) beforeRowQuery(scope *gorm.Scope) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.Scope)  { c.after(scope, "") }

func (c *callbacks) before(scope *gorm.Scope) {
	val, ok := scope.Get(parentSpanGormKey)
	if !ok {
		return
	}
	parentSpan := val.(opentracing.Span)
	tr := parentSpan.Tracer()
	sp := tr.StartSpan("sql", opentracing.ChildOf(parentSpan.Context()))
	ext.DBType.Set(sp, "sql")
	scope.Set(spanGormKey, sp)
}

func (c *callbacks) after(scope *gorm.Scope, operation string) {
	val, ok := scope.Get(spanGormKey)
	if !ok {
		return
	}
	sp := val.(opentracing.Span)
	if operation == "" {
		operation = strings.ToUpper(strings.Split(scope.SQL, " ")[0])
	}
	ext.Error.Set(sp, scope.HasError())
	ext.DBStatement.Set(sp, scope.SQL)
	sp.SetTag("db.table", scope.TableName())
	sp.SetTag("db.method", operation)
	sp.SetTag("db.err", scope.HasError())
	sp.SetTag("db.vars", scope.SQLVars)
	sp.SetTag("db.count", scope.DB().RowsAffected)
	sp.Finish()
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	// gorm does some magic, if you pass CallbackProcessor here - nothing works
	switch name {
	case "create":
		db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		db.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		db.Callback().RowQuery().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}
func Close() {
	_ = db_write.Close()
}
