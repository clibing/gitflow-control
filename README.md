### 快速创建分支

```bash
    git-feat
    git-fix
    git-docs
    git-style 
    git-refactor
    git-chore
    git-test
    git-hotfix
```

### 编译与安装

* 1.采用源码安装
  
  自行安装go开发环境。

  ```shell
  克隆代码到go的工作目录
  git clone https://github.com/clibing/gitflow-control $GOPATH/src/github.com/clibing/gitflow-control

  打开项目
  cd $GOPATH/src/github.com/clibing/gitflow-control 

  更新依赖
  go mod tidy

  编译
  make

  选择对应的平台安装, 以下以macos为例, 有的机器可能需要sudo需要输入管理员密码
  ./disk/gitflow-control_darwin_amd64 install 或者 sudo ./disk/gitflow-control_darwin_amd64 install 
  ```
* 2.下载已经编译的二进制编译
  ```
  下载github.com release的安装包自行安装
  curl -o gitflow-control 下载地址,查看github.com的release最新的下载地址， https://github.com/clibing/gitflow-control/releases

  chmod +x gitflow-control 
  gitflow-control install 或者 sudo gitflow-control install
  ```
* 3.安装与卸载
  ```
  安装：install 上面已经说了

  卸载：gitflow-control uninstall

  安装可以执行目录,例如`gitflow-control install -p /usr/local/bin`, 卸载时候同样需要执行目录 `gitflow-control uninstall -p /usr/local/bin`
  ```
* 4.安装生成的文件
  以macos为例
  ```
  ls -al
  /usr/local/bin/git-chore -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-ci -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-docs -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-feat -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-fix -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-hotfix -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-refactor -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-style -> /usr/local/bin/gitflow-control
  /usr/local/bin/git-test -> /usr/local/bin/gitflow-control
  ```

  ls -al 
  ```
  /Users/clibing/.gitflow-control
  ├── control.yaml
  └── hooks
    └── commit-msg -> /usr/local/bin/gitflow-control
  ```

  查看配置文件与说明
  ```
  mode: auto          # 模式，默认audo, 还有first(issue存在头部), standard 标准格式(见底部)
  issue:              # issue 的相关配置
    prefix-url: []    # 当mode为auto时，如果当前git仓库的remote的url以某一个为前缀，会开启first头部模式
    left-marker: ""   # 开启头部模式包装Footer内容的左标识符
    right-marker: ""  # 开启头部模式包装Footer内容的右标识符
    value: ""         # 会记录最近一次提交时使用的issue号。当重复提交代码时进行复用自动填充， 此值可以通过 `git issue -v ""` 设置
  ```

### 使用

* 1.快速创建分支
  ```
  git feat 分支的名字(新的功能)
  git fix 分支的名字(常规bug修复)
  git docs 分支的名字(文档相关)
  git style  分支的名字(格式（不影响代码运行的变动）)
  git refactor 分支的名字(重构（即不是新增功能，也不是修改bug的代码变动））
  git chore 分支的名字(构建过程或辅助工具的变动）
  git test 分支的名字(增加测试）
  git hotfix  分支的名字(紧急修复线上bug）
  ```
* 2.提交代码
  ```
  git add .
  git ci 后按照提示进行填写信息
  ```
  主要如果开启issue头部模式，需要在Footer项中直接填写issue的号即可
* 3.修改最近一次使用的issue号
  ```
  git issue -v "WEB-001"
  ```
  主要issue的格式是 "英文-数字"的格式

### 调研
* [git 自定义命令](https://blog.csdn.net/danpu0978/article/details/107276394)
* [Go 每日一库之 bubbletea](https://cloud.tencent.com/developer/article/1839581)
* [Angular community specification. Git commit 规范](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#heading=h.greljkmo14y0)
* [本项目主要的参考与复用的项目](https://github.com/mritd/gitflow-toolkit)