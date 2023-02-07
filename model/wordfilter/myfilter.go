// @Author Bing 
// @Date 2023/2/7 14:37:00 
// @Desc
package wordfilter


import (
	"github.com/importcjj/sensitive"
)
var Filter * sensitive.Filter

func init(){
	Filter = sensitive.New()
	Filter.LoadNetWordDict("https://raw.githubusercontent.com/importcjj/sensitive/master/dict/dict.txt")
}