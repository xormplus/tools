# xorm tools
## 说明
* xorm的辅助工具，本工具目前主要提供xorm的SqlMap配置文件和SqlTemplate模板批量加密功能。
* 目前支持AES，DES，3DES，RSA四种加密算法。其中AES，DES，3DES并非标准实现，有内置补足key，与[https://github.com/xormplus/xorm](https://github.com/xormplus/xorm) 库中的解密算法对应。
* 本工具使用Sciter的Golang绑定库 [sciter](https://github.com/oskca/sciter) 开发。由于主要是试用Sciter，所以逻辑相关的代码组织的并不是很规整，例如有些方法明显可以抽成接口方式。
* win64下可运行文件下载：[tools.2016.06.15.win64.zip](https://github.com/xormplus/tools/releases/download/v2016.06.15-alpha/tools.2016.06.15.win64.zip)，其他环境请自行编译

![](http://i.imgur.com/cR1lh2J.png)

