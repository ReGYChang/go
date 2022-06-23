- [Git](#git)
  - [Installation](#installation)
  - [Create Repository](#create-repository)
    - [git init](#git-init)
    - [git status](#git-status)
    - [git add](#git-add)
    - [git commit](#git-commit)
- [Version Control](#version-control)
    - [git log](#git-log)
    - [git diff](#git-diff)
  - [Rewriting History](#rewriting-history)
    - [git reset](#git-reset)
    - [git reflog](#git-reflog)

# Git

為了更方便地管理 Linux 程式碼, Linus 花了兩週自己用 C 開發了一套分散式版本控制系統, 即現在大家熟知的 Git

Linus 因為痛恨集中式的版本控制系統, 因此一直不願意使用 CVS 或 SVN 等工具

集中式的版本控制系統 codebase 是儲存在中央 server, 在做版控時要先從中央 server 取得最新的版本, 修改完程式碼再推送回去

集中式版控必須要透過網絡才能運作, 而分散式版控則沒有中央 server, 每個人 local 端都有一個完整的 codebase, 多人協同時只需要把各自的修改互相推送給對方即可看到對方的修改

## Installation

## Create Repository

Repository 中所有的文件都可以被 Git 管理, 舉凡每個文件的修改, 刪除都能被 Git 追蹤, 以便查看 codebase 歷史紀錄或是還原到某個時間點

### git init

以下範例創建一個 repo:

```shell
➜  mkdir learngit
➜  cd learngit
➜  pwd
/Users/regy/learngit
➜  learngit git init                                
hint: Using 'master' as the name for the initial branch. This default branch name
hint: is subject to change. To configure the initial branch name to use in all
hint: of your new repositories, which will suppress this warning, call:
hint: 
hint:   git config --global init.defaultBranch <name>
hint: 
hint: Names commonly chosen instead of 'master' are 'main', 'trunk' and
hint: 'development'. The just-created branch can be renamed via this command:
hint: 
hint:   git branch -m <name>
Initialized empty Git repository in /Users/regy/Github/test/learngit/.git/
```

如此一來一個空的 repo 就建好了, 目錄下產生了一個 `.git` 目錄, 其為 Git 用來跟蹤管理 repo, 不要隨意動到其中的文件, 以免破壞了 git repo

所有的版控系統只能追蹤文本文件的改動, 如 `.txt`, 程式碼等, 其會紀錄並顯示每次的文本改動, 如在第 5 行新增了一個單字 **Linux**, 在第 8 行刪除了一個單字 **mfer**, 針對 binary file 就無法追蹤其變化

> 建議使用標準的 `UTF-8` 編碼

下面示範如何將文件新增到版控系統中:

README.md

```md
Git is a version control system.
Git is free software.
```

首先將 `README.md` 放到 `learngit` 目錄下, 代表由此 git repo 來作管理

### git status

使用 `git status` 查看當前 git repo 的狀態:

```shell
➜  learngit git:(master) ✗ git status
On branch master

No commits yet

Untracked files:
  (use "git add <file>..." to include in what will be committed)
        README.md

nothing added to commit but untracked files present (use "git add" to track)
```

### git add

此時 `README.md` 雖然被 Git 偵測到, 但目前屬於 `Untracked files`, 表示尚未是 Git 追蹤的對象, 需要使用 `git add` 將文件新增到 git stagin area 中才能將 `README.md` 加入到追蹤對象:

```shell
➜  learngit git:(master) ✗ git add README.md
```

再使用 `git status` 查看會發現此時狀態會從 `Untracked files` 變成 `Changes to be committed`, 表示放在索引中的文件即將會被提交成一個新版本(commit)

```shell
➜  learngit git:(master) ✗ git status       
On branch master

No commits yet

Changes to be committed:
  (use "git rm --cached <file>..." to unstage)
        new file:   README.md
```

### git commit

此時可以用 `git commit` 提交一個新版本到 git repo:

```shell
➜  learngit git:(master) ✗ git commit -m 'first commit'                 
[master (root-commit) cbb0c14] first commit
 1 file changed, 2 insertions(+)
 create mode 100644 README.md
```

最後使用 `git status` 可以看到以下訊息:

```shell
➜  learngit git:(master) git status                   
On branch master
nothing to commit, working tree clean
```

表示已將 `README.md` 提交成一個 commit, 所以目前工作目錄上已經清空了

# Version Control

當修改完程式碼之後透過 `git commit` 提交到 git repo 中, 如果哪天程式碼被改壞了或是誤刪了什麼文件, 依然可以從任意 commit 恢復而不會造成無法彌補的傷痛

再嘗試修改文件, 並將修改提交到 git repo(修改 README.md 如下):

```
Git is a distributed version control system.
Git is free software distributed under the GPL.
```

此時用 `git status` 查看:

```shell
➜  learngit git:(master) git status       
On branch master
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   README.md

no changes added to commit (use "git add" and/or "git commit -a")
```

狀態又變為 `Changes not staged for commit` 並顯示 `README.md` 文件被修改, 然後嘗試提交:

```shell
➜  learngit git:(master) ✗ git commit -am 'second commit' 
[master 5dcfc65] second commit
 1 file changed, 3 insertions(+), 2 deletions(-)
```

注意這裡 `git commit -am` 中的 `a` 參數等價於 `git add`

### git log

接著使用 `git log` 可以查看 git repo 中的 history commit:

```shell
commit 5dcfc65acad6776f00c9375648ccb8b83315e603 (HEAD -> master)
Author: ReGYChang <p714140432@gmail.com>
Date:   Wed Jun 22 22:24:24 2022 +0800

    second commit

commit cbb0c143579ff7d2c21cd8c66d00d2a02458ae64
Author: ReGYChang <p714140432@gmail.com>
Date:   Wed Jun 22 22:08:24 2022 +0800

    first commit
```

Git 的 `commit id` 由一串雜湊值表示, 這是一個 `SHA1` 計算出來的一個數字, 以十六進制表示, git 就是透過 `commit id` 來實現版本控制

若覺得 `git log` 輸出內容太多, 也可以加上參數 `--oneline`:

```shell
5dcfc65 (HEAD -> master) second commit
cbb0c14 first commit
```

### git diff

如果要比較文件與上個版本的差異, 可以使用 `git diff` 查看:

```shell
➜  learngit git:(master) git diff 5dcfc65 cbb0c14

diff --git a/README.md b/README.md
index ce32b56..d8036c1 100644
--- a/README.md
+++ b/README.md
@@ -1,3 +1,2 @@
-```
-Git is a distributed version control system.
-Git is free software distributed under the GPL.
\ No newline at end of file
+Git is a version control system.
+Git is free software.
\ No newline at end of file
```

## Rewriting History

### git reset

Git 中 `HEAD` 表示當前版本, 如果要回退到上個版本 `first commit`, 可以使用 `git reset`:

```shell
➜  learngit git:(master) git reset --hard cbb0c14                                          
HEAD is now at cbb0c14 first commit
```

再用 `git log` 指令可以發現, 此時 `README.md` 文件的版本已經回退到 `first commit` 的版本, 而剛剛最新的版本 `second commit` 已經不見了:

```shell
commit cbb0c143579ff7d2c21cd8c66d00d2a02458ae64 (HEAD -> master)
Author: ReGYChang <p714140432@gmail.com>
Date:   Wed Jun 22 22:08:24 2022 +0800

    first commit
```

若想再回到 `second commit` 的版本, 就再使用一次 `git reset` 即可:

```shell
➜  learngit git:(master) git reset --hard 5dcfc65       
HEAD is now at 5dcfc65 second commit
```

Git 版本回退的速度非常快, 在內部有個指向當前版本的 `HEAD` pointer, 當回退版本時 git 只是把 `HEAD` pointer 從指向 `second commit` 改成指向 `first commit` 並把工作區文件更新了:

```
┌────┐
│HEAD│
└────┘
   │
   └──> ○ append GPL
        │
        ○ add distributed
        │
        ○ wrote a readme file

👇

┌────┐
│HEAD│
└────┘
   │
   │    ○ append GPL
   │    │
   └──> ○ add distributed
        │
        ○ wrote a readme file
```

那如果回退到某個版本後後悔, 想恢復到新版本怎麼辦卻找不到新版本的 `commit id` 怎麼辦?

### git reflog

在 git 中總是有後悔藥可以吃, 可以使用 `git reflog` 來查看之前使用過的指令:

```shell
5dcfc65 (HEAD -> master) HEAD@{0}: reset: moving to 5dcfc65
cbb0c14 HEAD@{1}: reset: moving to cbb0c14
5dcfc65 (HEAD -> master) HEAD@{2}: commit: second commit
cbb0c14 HEAD@{3}: commit (initial): first commit
```

就可以找到更新版本的 `commit id` 並使用 `git reset` 來移動 `HEAD` pointer