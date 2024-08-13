##  brutepool：一旦某个子goroutine成功，就停止全部goroutine。

### 简介
* 目前我看到的goroutine池都没有提前终止功能，就意味着字典有多大就要跑多久。
* 于是参考HackPool，有了brutepool。
* 一旦字典里的某个字符串成功爆破了，就取消所有的goroutine，结束字典内后面数据的爆破。然后运行回调函数。

### 使用说明

1. 先定义一个类型为[]interface{}的字典(BruteList)。

    ```
        passList := ['123456','admin','111111','root']
        bruteList := make([]interface{}, len(passList), len(passList))
        for i := range passList {
            bruteList[i] = passList[i]
        }
    ```
   
2. 再定义一个返回bool的爆破函数(BruteFunc)。
    ```
    func bruteFunc(passwd interface{}) bool {
            if passwd成功{
                return true
                }
            return false
       }
    ```
3. 导入github.com/hanbufei/brutepool包。
   ```
   import github.com/hanbufei/brutepool
   ```

4.　调用brutepool的New和Run即可（默认线程数4，回调函数是打印成功的值）。如果不想使用默认配置，则使用第5步代替。

```
    p := brutepool.New(bruteList, bruteFunc)
    p.Run()
```
5. 直接定义brutepool并Run。
```
    p := &brutepool.BrutePool{
        BruteList:       bruteList,
        Concurrency:     你定义的线程数,
        BruteFunc:       bruteFunc,
        SuccessCallBack: 你定义的回调函数,
        queues:          make(chan interface{}),
    }
    p.Run()
```

### 示例
参考 cmd/main.go

### 对比测试
![img.png](img.png)