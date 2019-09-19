# go 语言web框架
gin只是封装请求上下文，数据库连接、日志等功能都没有，这里封装了一下方便后面开箱即用

# 这里有几个疑问点

为什么代码放开internal里面，而没有放在pkg里面?\
internal里面代码意思是不共享不能被别人的项目import的，如果是开源组建的话放在pkg里面的包可以被别人引用，例如我们import的gin都是在pkg下面的。这里只是一个web框架不用做开源引用，所以怎么写都行

主从多个数据库连接怎么维护? \
`internal/db/mysql`下面封装了NewDB方法，照例在来一个全局变量就好了

日志怎么维护的? \
日志使用的是第三方扩展，logger/zap封装了一下啊zap的功能，如果想换一个扩展或者想自己实现一下这写方法就好了

路由怎么用? \
路由的话就完全使用gin的路由

mvc怎么实现? \
web mvc太经典了，不是说换种语言就不用mvc了，不过现在v应该是很少用到了，gin也是支持模版渲染，不多说，m的话用的是gorm，c的话我个人是用handler来代替柔和，作用是验证参数、柔和方法完成当前功能

配置文件? \
配置的话，这里也是用全局变量来维护的，internal/config.Conf 结构体保存了变量配置