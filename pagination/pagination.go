// 分页控制器

package pagination

import (
	"math"
	"strconv"
)

/*
	//示例代码
	var params pager.Params
	params.TotalPage = 20
	params.URIPrefix = "/card/list?id=1&page="
	params.CurrentPage = 20
	//翻页样式
	page := pager.Turn(&params)
	log.Println("首页：", page.First)
	log.Println("上页：", page.Prev)
	log.Println("下页：", page.Next)
	log.Println("尾页：", page.End)
	//页码样式
	params.ShowNumber = 5
	page2 := pager.NumberList(&params)
	for _, v := range page2 {
		log.Println("页码", params.URIPrefix + strconv.Itoa(v))
	}
*/

// Params 输入参数
type Params struct {
	URIPrefix   string // URL前缀
	TotalPage   int    // 总页数
	CurrentPage int    // 当前页
	ShowNumber  int    // 显示页码的个数
}

// ControllerParams 控制器样式的变量
type ControllerParams struct {
	First string // 首页
	Prev  string // 上页
	Next  string // 下页
	End   string // 尾页
}

// 控制器样式
func Controller(params *Params) *ControllerParams {
	var result ControllerParams
	if params.TotalPage <= 1 {
		return &result
	}

	if params.CurrentPage > 1 {
		result.First = params.URIPrefix + "1"
		result.Prev = params.URIPrefix + strconv.Itoa(params.CurrentPage-1)
	}
	if params.CurrentPage < params.TotalPage {
		result.Next = params.URIPrefix + strconv.Itoa(params.CurrentPage+1)
		result.End = params.URIPrefix + strconv.Itoa(params.TotalPage)
	}
	return &result
}

// 页码列表样式
func NumberList(params *Params) []int {
	if params.CurrentPage > params.TotalPage {
		params.CurrentPage = params.TotalPage
	}
	if params.CurrentPage <= 0 {
		params.CurrentPage = 1
	}
	var pages []int
	switch {
	case params.CurrentPage >= params.TotalPage-params.ShowNumber && params.TotalPage > params.ShowNumber:
		start := params.TotalPage - params.ShowNumber + 1
		pages = make([]int, params.ShowNumber)
		for i := range pages {
			pages[i] = start + i
		}
	case params.CurrentPage >= 3 && params.TotalPage > params.ShowNumber:
		start := params.CurrentPage/2 + 1
		pages = make([]int, params.ShowNumber)
		for i := range pages {
			pages[i] = start + i
		}
	default:
		pages = make([]int, int(math.Min(float64(params.ShowNumber), float64(params.TotalPage))))
		for i := range pages {
			pages[i] = i + 1
		}
	}
	return pages
}
