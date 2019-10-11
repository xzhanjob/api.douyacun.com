github master push 自动触发部署到博客

# 文章配置 
index: /:topic/_doc/:id
```yaml
Title: 文章标题
# 注意用 `,` 分隔
Keywords: 关键字,seo优化使用
Description: 文章
Author: douyacun
Date: 2019-09-19 18:03:32
LastEditTime: 2019-10-09 14:36:06
typora-root-url: ./assert
```
# ci 配置
- travis ci 自动部署
- ci部署完成执行go脚本
- douyacun.yml
```yaml
# 作者
Author: douyacun
# 作者邮箱 (默认当作全局配置，优先展示文章配置的邮箱）
Email: douyacun@gmail.com
# 作者github连接
Gihutb: https://github.com/douyacun
# 微信公众号
WeChatSubscription: douyacun
# 微信公众号二维码
WeChatSubscriptionQrcode: /assert/douyacun_qrcode.jpg
Topics: # 话题
  # golang话题
  Golang:
    # 建议icon使用svg
    Icon: 话题icon
    # 文章所在目录 
    Dir: /go
    # 图片路径
    Assert: /go/assert
    # 文章目录支持1级，后续考虑多级
    Articles:
      # 文章外部图片建议在内部进行配置
      - 函数方法接口.md
      - 数组切片引用.md
      - json解析技巧.md
  # redis话题
  redis:
```
- go脚本解析douyacun.yml
- 根据douyacun.yml中配置的文章路径收集文章
- 文章分析 写入 elstaticsearch

# 订阅最新消息
- index: subscriber
- source: 
```json
{
  "email": "douyacun@gmail.com",
  "date": "2019-10-10 17:52:32"
}
```

# todo
- [ ] markdown 本地跳转
- [ ] 文章关键词提取
- [ ] travis ci 自动部署
- [x] 图片提取
- [ ] git提取文件创建时间
- [ ] markdown 视频