# goip #

调用ipip.net的免费版IP数据库查询IP所在地。根据ipip.net提供的php版代码硬翻。


### 使用 ###

```go
package main

import (
	"fmt"

	"github.com/jiazhoulvke/goip"
)

func main() {
	if err := goip.SetDBPath("ip数据库路径"); err != nil {
		panic(err)
	}
	location, err := goip.Find("8.8.8.8")
	if err != nil {
		panic(err)
	}
	fmt.Println(location.Country)  //国家
	fmt.Println(location.Province) //省份
	fmt.Println(location.City)     //城市
	fmt.Println(location.District) //区县
}
```

