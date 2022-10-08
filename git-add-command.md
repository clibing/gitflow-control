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