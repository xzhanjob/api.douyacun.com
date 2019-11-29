github master push 自动触发部署到博客

[![Build Status](https://travis-ci.org/douyacun/api.douyacun.com.svg?branch=master)](https://travis-ci.org/douyacun/api.douyacun.com)

# 全局配置
**douyacun.yml 全局配置**
```yaml
# 作者
Author: douyacun
# 作者邮箱 (默认当作全局配置，优先展示文章配置的邮箱）
Email: douyacun@gmail.com
# 密钥
key: 1231345
# 作者github连接
Gihutb: https://github.com/douyacun
# 微信公众号
WeChatSubscription: douyacun
# 微信公众号二维码
WeChatSubscriptionQrcode: /assert/douyacun_qrcode.jpg
Topics: # 话题
  # golang话题
  Golang:
      # 文章外部图片建议在内部进行配置
      - 函数方法接口.md
      - 数组切片引用.md
      - json解析技巧.md
```
# 文章配置 
index: /:topic/:id
```yaml
# 如果读取不到标题，使用文件名作为标题
Title: 文章标题
# 注意用 `,` 分隔
Keywords: 关键字,seo优化使用
Description: 文章
# 没有使用下面全局配置中的author
Author: douyacun
# date：文件创建时间，date > git首次提交时间，默认会取git版本中的首次提交日期 
Date: 2019-09-19 18:03:32
# LastEditTime: 文件更新日期，LastEditTime < git版本最后一次提交日期，默认会取git版本最后提交日期
LastEditTime: 2019-10-09 14:36:06
```
- 文章唯一性标识，md5(douyacun.yml->key, "+", "文章所属话题(Golang)", "+", "文件名称")

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
- [ ] markdown 本地跳转, 1-go-cannel.md
- [ ] 封面功能，配置没有文件，取文档第一张图片作为封面
- [ ] 图片压缩功能
- [x] git提取文件创建时间
- [x] 文章关键词提取
- [x] travis ci 自动部署
- [x] 图片提取
- [x] 文章封面
- [x] 首页文章分页, 文章按更新时间排序
- [x] 前端js文件404问题(nextjs link 默认会是预加载页面： <Link> will automatically prefetch pages in the background )
- [x] 制作favicon.ico
- [x] 适配手机端
- [x] 增加友盟
