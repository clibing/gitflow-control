git 自定义命令
https://blog.csdn.net/danpu0978/article/details/107276394

例如

系统内置的command
```
git branch
```

增加一个新的command me
```
git me 输入当前$USER
```

touch /usr/local/bin/git-me
```sh
#!/bin/bash

echo $USER
```

chmod +x /usr/local/bin/git-me 

git me 输入当前用户名s


增加常用的命令

快速创建分支
```bash

    git-feat
    git-fix
    git-docs
    git-style 
    git-refactor
    git-chore
    git-test
    git-hotfix
    git-perf
```
清理issue的记忆
```bash
    git-issue [init|add|remove|reset] 可以增加多个 在ci时可以选择对应的issue
    git issue init -prefix [ -suffix ] -show header
    show: 代表issue的展示位置, 
        允许header: 如果在头，不进行Closes标记
        允许footer: 默认以Closes 开头: #001, #002s
```
提交
```bash
    git-ci
```

调研
* [https://cloud.tencent.com/developer/article/1839581](https://cloud.tencent.com/developer/article/1839581)