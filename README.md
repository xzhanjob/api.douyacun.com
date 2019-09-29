github master push 自动触发部署到博客

# elstaticsearch 存储文档
index: /:topic/_doc/:id
```json
{
    "title": "标题",
    "keywords": "关键字",
    "Description": "文件描述",
    "author": "作者",
    "content": "这里存储的是文章内容"
}
```
- 支持中文分词
- 搜索高亮

# CI
- travis ci 自动部署
- ci部署完成执行go脚本
- douyacun.yml
```yaml
# 作者
author: douyacun
# 作者邮箱 (默认当作全局配置，优先展示文章配置的邮箱）
email: douyacun@gmail.com
# 作者github连接
gihutb: https://github.com/douyacun
# 微信公众号
WeChatSubscription: douyacun
# 微信公众号二维码
WeChatSubscriptionQrcode: /assert/douyacun_qrcode.jpg
topics: # 话题
  # golang话题
  golang:
    # 建议icon使用svg
    icon: 话题icon
    # 文章目录支持1级，后续考虑多级
    article:
      # 文章外部图片建议在内部进行配置
      - json解析技巧.md
      - select.md
  # redis话题
  redis:
  mysql:
```
- go脚本解析douyacun.yml
- 根据douyacun.yml中配置的文章路径收集文章
- 文章分析 写入 elstaticsearch

# website

route:
- / 网站根目录
- /posts/:id 文章详情页
- /topic/:id 话题