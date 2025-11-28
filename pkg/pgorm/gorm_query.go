package pgorm

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QueryList_ 通用列表查询接口
// Deprecated
// e.g. 在 db/models层调用
//	var total int64
//	var err error
//	var students []Student
//	total, err = pgorm.QueryList(c, DB(), []string{"模糊查询字段定义"},
//		map[string]interface{}{
//			"精确字段定义": "",
//		}, &students)
//	pgin.CheckErr(err)
//
//	pgin.OK(c, gin.H{
//		"count": total,
//		"info": students,
//	})
func QueryList_(c *gin.Context, db *gorm.DB, searchFields []string,
	filterFields map[string]interface{}, models interface{}) (total int64, err error) {
	page := queryIntDefault(c, "page", 1)
	pageSize := queryIntDefault(c, "page_size", 10)
	search := c.Query("search")
	order := c.Query("order")
	filter := make(map[string]interface{})
	for k, v := range filterFields {
		if s := c.Query(k); s != "" {
			value := reflect.ValueOf(v)
			switch value.Kind() {
			case reflect.Int:
				if tmp, err := strconv.Atoi(s); err == nil {
					filter[k] = tmp
				}
			case reflect.Bool:
				tmp := !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
				filter[k] = tmp
			default:
				// string
				filter[k] = s
			}
		}
	}
	q := &PaginationQuery{
		Page: page,
		PageSize: pageSize,
		Search: search,
		SearchFields: searchFields,
		Order: order,
		Filter: filter,

		db: db,
		model: models,
	}

	return q.Query(models)
}

func queryIntDefault(c *gin.Context, key string, d int) int {
	q := c.Query(key)
	if q == "" {
		return d
	}
	if v, err := strconv.Atoi(q); err == nil {
		return v
	}
	return d
}

func queryBoolDefault(c *gin.Context, key string, d bool) bool {
	q := c.Query(key)
	if q == "" {
		return d
	}
	s := strings.ToLower(strings.TrimSpace(q))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
}

/********************************************************************************************/

type queryPara struct {
	Page 		int	   `form:"page"`
	PageSize	int	   `form:"page_size"`
	Search 		string `form:"search"`
	Order 		string `form:"order"`
	// 精确匹配优先使用 filter_map
	FilterMap   map[string]interface{} `form:"filter_map"`
	// 其他过滤字段直接透传
	FilterOr   		map[string]interface{}	`form:"filter_or_map"`
	NotFilter   	map[string]interface{}	`form:"not_filter_map"`
	GreaterThan 	map[string]interface{}	`form:"gte_map"`
	LessThan 		map[string]interface{}	`form:"lte_map"`
}

// QueryList 通用列表查询接口
// e.g. 在 db/models层调用
//	var total int64
//	var err error
//	var students []Student
//	total, err = pgorm.QueryList(c, DB(), []string{"模糊查询字段定义"},
//		map[string]interface{}{
//			"精确字段定义": "",
//		}, &students)
//	pgin.CheckErr(err)
//
//	pgin.OK(c, gin.H{
//		"count": total,
//		"info": students,
//	})
func QueryList(c *gin.Context, db *gorm.DB, searchFields []string,
	filterFields map[string]interface{}, models interface{},
	value ...interface{}) (total int64, err error) {
	q := GetRequestQuery(c, searchFields, filterFields, value...)
	q.db = db
	q.model = models

	return q.Query(models)
}

// GetRequestQuery 通用获取请求查询条件
func GetRequestQuery(c *gin.Context, searchFields []string, filterFields map[string]interface{},
	value ...interface{}) *PaginationQuery {
	para := queryPara{
		Page: 1, PageSize: 10,
	}
	c.ShouldBindQuery(&para)
	var transFields map[string]string
	var sortFields []string
	if len(value) > 0 {
		if v0, ok := value[0].(map[string]string); ok {
			// 针对精确匹配/ 适配db到前端字段名称映射不一致的场景
			transFields = v0
		}
		if len(value) > 1 {
			if v1, ok := value[1].([]string); ok {
				// 用于限制排序字段的范围/ 或者仅仅是用于防止不存在的字段报错
				sortFields = v1
			}
		}
	}

	filter := make(map[string]interface{})
	for k, v := range filterFields {
		q := k
		if transFields != nil {
			if k2, ok := transFields[k]; ok { q = k2 }
		}

		str := ""
		if len(para.FilterMap) > 0 {
			if fv, ok := para.FilterMap[q]; ok {
				if str, ok = fv.(string); !ok {
					b, _ := json.Marshal(fv)
					str = string(b)
				}
			}
		} else {
			str = c.Query(q)
		}
		if str != "" {
			val := reflect.ValueOf(v)
			switch val.Kind() {
			case reflect.Int:
				var arr []int64
				if err := json.Unmarshal([]byte(str), &arr); err == nil {
					// 整型数组支持
					filter[k] = arr
					continue
				}
				if tmp, err := strconv.Atoi(str); err == nil {
					filter[k] = tmp
				}
			case reflect.Bool:
				tmp := !(str == "" || str == "0" || str == "no" || str == "false" || str == "none")
				filter[k] = tmp
			default:
				// string
				var arr []string
				if err := json.Unmarshal([]byte(str), &arr); err == nil {
					// 字符串数组支持
					filter[k] = arr
					continue
				}
				filter[k] = str
			}
		}
	}

	order := para.Order
	if order != "" && len(transFields) > 0 {
		transFields2 := map[string]string{}
		for k, v := range transFields {
			transFields2[v] = k
		}
		fields := strings.Split(strings.ReplaceAll(strings.ReplaceAll(order, "+", ""), "-", ""), ",")
		for _, v := range fields {
			if v2, ok := transFields2[v]; ok {
				order = strings.ReplaceAll(order, v, v2)
			}
		}
	}
	if order != "" && len(sortFields) > 0 {
		fields := strings.Split(strings.ReplaceAll(strings.ReplaceAll(order, "+", ""), "-", ""), ",")
		for _, v := range fields {
			exist := false
			for _, v2 := range sortFields {
				if v == v2 {
					exist = true
					break
				}
			}
			if !exist {
				order = ""
				break
			}
		}
	}
	if order == "" { order = "-updated_at" }

	q := &PaginationQuery{
		Page: para.Page,
		PageSize: para.PageSize,
		Search: para.Search,
		SearchFields: searchFields,
		Order: order,
		Filter: filter,
		FilterOr: para.FilterOr, NotFilter: para.NotFilter,
		GreaterThan: para.GreaterThan, LessThan: para.LessThan,
	}
	if q.FilterOr == nil { q.FilterOr = map[string]interface{}{} }
	if q.NotFilter == nil { q.NotFilter = map[string]interface{}{} }
	if q.GreaterThan == nil { q.GreaterThan = map[string]interface{}{} }
	if q.LessThan == nil { q.LessThan = map[string]interface{}{} }
	return q
}
