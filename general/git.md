- [Git](#git)
  - [Installation](#installation)
  - [Create Repository](#create-repository)

# Git

為了更方便地管理 Linux 程式碼, Linus 花了兩週自己用 C 開發了一套分散式版本控制系統, 即現在大家熟知的 Git

Linus 因為痛恨集中式的版本控制系統, 因此一直不願意使用 CVS 或 SVN 等工具

集中式的版本控制系統 codebase 是儲存在中央 server, 在做版控時要先從中央 server 取得最新的版本, 修改完程式碼再推送回去

集中式版控必須要透過網絡才能運作, 而分散式版控則沒有中央 server, 每個人 local 端都有一個完整的 codebase, 多人協同時只需要把各自的修改互相推送給對方即可看到對方的修改

## Installation

## Create Repository

Repository 中所有的文件都可以被 Git 管理, 舉凡每個文件的修改, 刪除都能被 Git 追蹤, 以便查看 codebase 歷史紀錄或是還原到某個時間點

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