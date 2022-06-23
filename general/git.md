- [Git](#git)
- [Installation](#installation)
- [Create Repository](#create-repository)
  - [git init](#git-init)
  - [git status](#git-status)
  - [git add](#git-add)
  - [git commit](#git-commit)
- [Working Directory, Staging Area & Repository](#working-directory-staging-area--repository)
  - [Working Directory](#working-directory)
  - [Staging Area](#staging-area)
  - [Repository](#repository)
- [Version Control](#version-control)
  - [git log](#git-log)
  - [git diff](#git-diff)
  - [git reset](#git-reset)
  - [git reflog](#git-reflog)
- [Undoing Changes](#undoing-changes)
- [Branch Management](#branch-management)
  - [Create & Merge Branch](#create--merge-branch)

# Git

為了更方便地管理 Linux 程式碼, Linus 花了兩週自己用 C 開發了一套分散式版本控制系統, 即現在大家熟知的 Git

Linus 因為痛恨集中式的版本控制系統, 因此一直不願意使用 CVS 或 SVN 等工具

集中式的版本控制系統 codebase 是儲存在中央 server, 在做版控時要先從中央 server 取得最新的版本, 修改完程式碼再推送回去

集中式版控必須要透過網絡才能運作, 而分散式版控則沒有中央 server, 每個人 local 端都有一個完整的 codebase, 多人協同時只需要把各自的修改互相推送給對方即可看到對方的修改

# Installation

# Create Repository

Repository 中所有的文件都可以被 Git 管理, 舉凡每個文件的修改, 刪除都能被 Git 追蹤, 以便查看 codebase 歷史紀錄或是還原到某個時間點

## git init

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

## git status

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

## git add

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

## git commit

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

# Working Directory, Staging Area & Repository

在 Git 中主要可以分為 `Working Directory`, `Staging Area` 和 `Repository` 三個區塊, 前面提到使用 `git add` 指令將檔案新增進 `Staging Area`(or index), 再使用 `git commit` 將 `Staging Area` 中的內容一次性移往 `Repository`

![working_directory](img/working_directory.png)

## Working Directory

> git status: `Untracked files`, `Changes not staged for commit`

`Working Deirectory` 即是我們在系統中看到的目錄, 比如 `/Users/regy/Github/test/learngit` 文件夾就是一個 `Working Directory`

## Staging Area

> git status: `Changes to be committed`, `new file`

又常被稱作 `index`, `Staging Area` 紀錄有哪些檔案即將要被提交到下一個 commit 版本中, 即要提交一個版本到 `Repository` 前必須要先更新 index status, 有變更才能提交成功

## Repository

> git status: `Committed`

前面提到`Working Directory` 有一個隱藏目錄 `.git`, 其中就包含了 Git 的 `Repository` 和 `Staging Area`

`Repository` 主要用來保存檔案或是程式碼, `Staging Area` 的資料提交到 `Repository` 後可以永久保存, 儲存相關內容的歷史修改紀錄及變更內容等

Git 跟蹤管理的是修改, 而非檔案, 如增加了一行, 刪除了一行, 更改某些字符, 甚至創建一個新的檔案都算是一個修改

而 `git commit` 只會將 `Staging Area` 中的修改提交到 `Repository` 中, 若修改沒有使用 `git add` 新增到 `Staging Area` 則不會被 `git commit` 提交到 `Repository`

如果覺得要先 `add` 再 `commit` 有點繁瑣, 也可以使用 `git commit -am`, 加上 `a` 參數的話即使沒有 `add` 也可以完成 `commit`, 但要注意的是 **`-a` 參數只對已經存在 `Repository` 中的檔案有效**, 若還沒新增進 `Respository` 的檔案(Untracked file) 也是無法提交成功

> 為什麼需要切分這麼多區域或階段呢?

原因是擁有 `Staging Area` 可以在操作上有更多的彈性與靈活性, 在 `commit` 到 `Repository` 之前可以針對不同狀況操作不同指令來控制檔案:

- 修改了三個檔案, 其中一個不想提交, 如何操作?
- 原先修改的檔案想放棄, 如何回到原來的版本?
- 尚未完成的檔案, 想先儲存可以怎麼做?
- 發現忘記切換分支怎麼辦, 想切回正確的分支?

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

## git log

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

## git diff

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

## git reset

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

## git reflog

在 git 中總是有後悔藥可以吃, 可以使用 `git reflog` 來查看之前使用過的指令:

```shell
5dcfc65 (HEAD -> master) HEAD@{0}: reset: moving to 5dcfc65
cbb0c14 HEAD@{1}: reset: moving to cbb0c14
5dcfc65 (HEAD -> master) HEAD@{2}: commit: second commit
cbb0c14 HEAD@{3}: commit (initial): first commit
```

就可以找到更新版本的 `commit id` 並使用 `git reset` 來移動 `HEAD` pointer


# Undoing Changes

如果不小心在 `README.md` 中加了一行:

```shell
➜  learngit git:(master) cat README.md                 

Git is a distributed version control system.
Git is free software distributed under the GPL.
Stupid boss mfer.%    
```

準備要提交時突然看到這行程式碼不能被老闆發現, 這時可以先用 `git status` 查看:

```shell
➜  learngit git:(master) ✗ git status
On branch master
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   README.md

no changes added to commit (use "git add" and/or "git commit -a")
```

這裡顯示 `Changes not staged for commit`, 代表這段修改還沒有被 `add` 到 `Staging Area`, 並提示可使用 `git restore` 來撤銷 `Working Directory` 中的變更:

```shell
➜  learngit git:(master) ✗ git restore README.md   
```

此時 `git restore` 捨棄掉在 `Working Directory` 中做的修改, 若修改已經 `add` 到 `Staging Area` 中, 則可以加上 `staged` 參數, 將在 `Staging Area` 的檔案修改回退到上一個狀態, 再使用一次 `git restore` 捨棄 `Working Directory` 中的修改:

```shell
➜  learngit git:(master) ✗ git add README.md        
➜  learngit git:(master) ✗ git status       
On branch master
Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
        modified:   README.md

➜  learngit git:(master) ✗ git restore --staged README.md   
➜  learngit git:(master) ✗ git restore README.md 
```

# Branch Management

前面說明了關於 Git 的基本使用方法, 再來就要介紹 Git 中最重要的 branch 功能

在開發軟體時可能同時會有多人在開發同一個功能或針對某個模組修復錯誤, 也可能會有多個發佈版本的存在且需要針對各個版本進行維護, Git 利用 branch 來支援這些功能

Branch 為了將版本修改紀錄的整體流程分開儲存, 讓切開的 branch 不受其他 branch 的影響, 所以在同一個 `Repository` 下可以同時進行多個不同版本的修改

Branch 也可以與其他的 branch merge, 如開發某個新功能的 branch 在開發完後再 merge 回 `main` branch, 如此一來能保證程式碼安全又不會因多人協同提交而互相影響

![branch](img/branch.png)

上面的圖顯示 `Repository` 中的三條 branch, `main` 一般來說代表在正式環境運行的程式碼版本, 而切分出來的 `Little Feature` branch 及 `Big Feature` branch 與 `main` branch 彼此互相獨立更新並記錄, 不止能讓三個不同版本的程式碼平行工作, 也能避免一些還在開發中的程式碼進到 `main` branch 影響 production 環境運作

當 `Repository` 第一次 commit 時, Git 會自動創建 `main` branch, 之後的 commit 在切換 branch 之前都會在 `main` branch 中做 commit

## Create & Merge Branch

Git 會將每次的 commit 串成一條 timeline, 這條 timeline 即為一個 branch, Git 默認會創建一條 `master` branch, 而 `HEAD` pointer 嚴格來說不是指向 commit, 而是指向 `master`, `master` 再指向 commit, 所以 `HEAD` 指向的就是當前所在的 branch

```shell
5dcfc65 (HEAD -> master) second commit
```

剛開始時 `master` branch 是一條 timeline, Git 用 `master` 指向最新的 commit, 再用 `HEAD` 指向 `master`, 就能確認當前的 branch 以及當前 branch 的 commit point:

```shell
➜  learngit git:(master) git branch
* master
```

![create_branch_1](img/create_branch_1.png)

每次 commit `master` branch 都會新增一個節點, 隨著不斷 commit `master` branch 也越來越長

當創建新的 branch 如 `dev` 時, Git 創建了一個新的 pointer `dev`, 其指向 `master` 相同的 commit point, 將 `HEAD` 指向 `dev` 即表示當前 branch 在 `dev` 上:

```shell
➜  learngit git:(master) git branch dev    
➜  learngit git:(master) git checkout dev                         
Switched to branch 'dev'
➜  learngit git:(dev) git branch      
* dev
  master
```

![create_branch_2](img/create_branch_2.png)

這邊可以觀察到 Git 創建一個 branch 的速度很快, 只需要增加一個 `dev` pointer, 並將 `HEAD` pointer 指向 `dev` 即可, `Working Directory` 的檔案不需做出任何變化

從現在起對 `Working Directory` 的修改和提交就是針對 `dev` branch 了, 如 commit 一次後 `dev` pointer 往前挪動一個節點, 而 `master` 則不動:

![create_branch_3](img/create_branch_3.png)

若在 `dev` 上的工作完成了, 也可以將 `dev` branch merge 到 `master` branch, 即將 `master` 指向 `dev` 的當前提交即可:

![create_branch_4](img/create_branch_4.png)

可以發現 Git merge branch 也只需要修改 pointer 指向, 而 `Working Directory` 內容也沒有改變

Merge branch 後甚至也可以刪除 `dev` branch, 即將 `dev` pointer 刪掉即可

![create_branch_5](img/create_branch_5.png)