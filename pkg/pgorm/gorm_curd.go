package pgorm

import (
	"fmt"
	"strings"

	"github.com/lemonkingstar/spider/pkg/plog"
	"gorm.io/gorm"
)

// PaginationQuery 分页查询
/*
- 分页查询：
?page=1&page_size=5
- 模糊查询：
?page=1&page_size=5&search=wuzi
- 单个字段升序
?page=1&page_size=5&order=name
- 单个字段降序
?page=1&page_size=5&order=-name
- 多个字段组合排序
?page=1&page_size=5&order=-name,id
- 字段精确匹配
?page=1&page_size=5&name=wuzi
- 字段多个值精确匹配
?page=1&page_size=5&name=["wuzi","liuzi"]
?page=1&page_size=5&id=[1,2]
- 多个字段精确匹配
?page=1&page_size=5&name=wuzi&id=1
- 字段精确匹配 - 优先匹配
?page=1&page_size=5&filter_map={"name":"wuzi"}
?page=1&page_size=5&filter_map={"id":1}
?page=1&page_size=5&filter_map={"id": [1,2]}
*/
type PaginationQuery struct {
	Page         int                    `form:"page" json:"page" cn:"获取页码"`
	PageSize     int                    `form:"page_size" json:"page_size" cn:"每页记录数"`
	Order        string                 `form:"order" json:"order" cn:"排序字段(e.g. -name)"`
	Filter       map[string]interface{} `form:"filter" json:"filter" cn:"精确查询 col-val 集合,支持 = 和 in"`
	FilterOr     map[string]interface{} `form:"filter_or" json:"filter_or" cn:"精确查询 col-val 集合,支持 = 和 in"`
	NotFilter    map[string]interface{} `form:"not_filter" json:"not_filter" cn:"精确查询 col-val 集合,支持 != 和 not in"`
	GreaterThan  map[string]interface{} `form:"greater_than" json:"greater_than" cn:"大于 col-val 集合,支持 >= "`
	LessThan     map[string]interface{} `form:"less_than" json:"less_than" cn:"小于 col-val 集合,支持 <="`
	Search       string                 `form:"search" json:"search" cn:"模糊匹配字符串,目前一次查询只支持一个模糊匹配的字符串"`
	SearchFields []string               `form:"search_fields" json:"search_fields" cn:"模糊查询匹配字段范围,数组存储所有需要模糊匹配的字段名"`

	db    *gorm.DB
	model interface{}
}

func (p *PaginationQuery) SetDB(db *gorm.DB) *PaginationQuery {
	p.db = db
	return p
}

func (p *PaginationQuery) SetModel(m interface{}) *PaginationQuery {
	p.model = m
	return p
}

func (p *PaginationQuery) Query(out interface{}) (total int64, err error) {
	p.db = p.db.Scopes(searchScope(p.SearchFields, p.Search))
	p.db = p.db.Scopes(filterScope(p.Filter)...)
	p.db = p.db.Scopes(filterOrScope(p.FilterOr))
	p.db = p.db.Scopes(notFilterScope(p.NotFilter)...)
	p.db = p.db.Scopes(greaterThanScope(p.GreaterThan)...)
	p.db = p.db.Scopes(lessThanScope(p.LessThan)...)
	p.db = p.db.Scopes(orderScope(p.Order))
	p.db.Model(p.model).Count(&total)
	err = p.db.Scopes(paging(p.Page, p.PageSize)).Find(out).Error
	return
}

// 分页
func paging(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	// fix
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = -1
		//return func(db *gorm.DB) *gorm.DB { return db.Limit(200) }
	}
	offset, limit := (page-1)*pageSize, pageSize
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

// 模糊搜索
func searchScope(fields []string, search string) (scope func(*gorm.DB) *gorm.DB) {
	if search == "" || len(fields) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	var queries []string
	var queryFields []interface{}

	for _, k := range fields {
		queries = append(queries, fmt.Sprintf("`%s` LIKE ?", k))
		queryFields = append(queryFields, fmt.Sprintf("%%%s%%", search))
	}
	scope = func(db *gorm.DB) *gorm.DB {
		if len(queries) > 0 {
			return db.Where(strings.Join(queries, " OR "), queryFields...)
		}
		return db
	}
	return
}

// 排序
func orderScope(order string) func(*gorm.DB) *gorm.DB {
	if order == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		var queries []string
		orders := strings.Split(order, ",")

		for _, k := range orders {
			if strings.HasPrefix(k, "-") {
				queries = append(queries, fmt.Sprintf("`%s` DESC", k[1:]))
			} else {
				queries = append(queries, k)
			}
		}
		return db.Order(strings.Join(queries, ", "))
	}
}

// 精确查询
func filterScope(filters map[string]interface{}) []func(*gorm.DB) *gorm.DB {
	if len(filters) == 0 {
		return nil
	}

	var scopes []func(*gorm.DB) *gorm.DB
	for k, v := range filters {
		field, value := k, v
		scope := func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s IN (?)", field), value)
		}
		scopes = append(scopes, scope)
	}
	return scopes
}

// Or条件查询
func filterOrScope(filters map[string]interface{}) func(*gorm.DB) *gorm.DB {
	if len(filters) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	var queries []string
	var queryFields []interface{}
	for field, value := range filters {
		queries = append(queries, fmt.Sprintf("`%s` IN (?)", field))
		queryFields = append(queryFields, value)
	}

	scope := func(db *gorm.DB) *gorm.DB {
		if len(queries) > 0 {
			return db.Where(strings.Join(queries, " OR "), queryFields...)
		}
		return db
	}
	return scope
}

// 精确查询不等于
func notFilterScope(filters map[string]interface{}) []func(*gorm.DB) *gorm.DB {
	if len(filters) == 0 {
		return nil
	}

	var scopes []func(*gorm.DB) *gorm.DB
	for k, v := range filters {
		field, value := k, v
		scope := func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` NOT IN (?)", field), value)
		}
		scopes = append(scopes, scope)
	}
	return scopes
}

// 查询大于某个值
func greaterThanScope(greaterThan map[string]interface{}) []func(*gorm.DB) *gorm.DB {
	if len(greaterThan) == 0 {
		return nil
	}

	var scopes []func(*gorm.DB) *gorm.DB
	for k, v := range greaterThan {
		field, value := k, v
		scope := func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` >= (?)", field), value)
		}
		scopes = append(scopes, scope)
	}
	return scopes
}

// 查询小于某个值
func lessThanScope(lessThan map[string]interface{}) []func(*gorm.DB) *gorm.DB {
	if len(lessThan) == 0 {
		return nil
	}

	var scopes []func(*gorm.DB) *gorm.DB
	for k, v := range lessThan {
		field, value := k, v
		scope := func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s` <= (?)", field), value)
		}
		scopes = append(scopes, scope)
	}
	return scopes
}

// WithTx 事务处理
func WithTx(db *gorm.DB, txFunc func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			plog.Errorf("[DB]transaction err %s", err.Error())
		}
		if err != nil {
			plog.Debugf("[DB]transaction rollback")
			if err2 := tx.Rollback().Error; err2 != nil {
				plog.Errorf("[DB]transaction rollback error %s", err2.Error())
			}
			return
		}
		plog.Debugf("[DB]transaction committed")
		if err = tx.Commit().Error; err != nil {
			plog.Errorf("[DB]transaction commit error %s", err.Error())
		}
	}()
	err = txFunc(tx)
	return
}
