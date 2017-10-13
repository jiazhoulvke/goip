# goip #

调用ipip.net的免费版IP数据库查询IP所在地。根据ipip.net提供的php版代码硬翻。


### 使用 ###

```go
if err := goip.SetDBPath("ip数据库路径");err!=nil {
  panic(err)
}
location, err:= goip.Find("8.8.8.8")
if err!=nil {
  panic(err)
}
fmt.Println(location.Country)
fmt.Println(location.Province)
fmt.Println(location.City)
fmt.Println(location.District)
```

