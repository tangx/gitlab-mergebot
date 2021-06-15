# gitlab-mergebot

gitlab merge request 助手

## Usage

### Mergebot Configuration

1. 选择任意一个公共账户， 或创建一个新账户。
2. 在**用户配置**中， 创建 `Access Token` ， 权限为 `api`。
3. 启动 `mergebot`

```bash
GITLAB_Endpoint="https://git.example.com"
GITLAB_PrivateToken="Adfasdfadfasdfaksdjf"
./gitlab-mergebot
```

![](./docs/mergebot/make-up.png)

### .mergebot.yml in Target Branch

1. 将 **mergebot** 加入 **项目** 或 **项目所在的 group**， 权限设置为 `Developer`

2. 在目标分支根目录下创建 `.mergebot.yml`， 如下

```yaml
# 最少 code review 人数
minReviewers: 2

# code review 有效名单
reviewers:
  - guanyu 
  - zhangfei
  - zhaoyun

# 合并人员名单
## 如果存在， mergebot 将在满足合并条件时，随机选择一位候选人并转让
## 如不存在， mergebot 将直接合并
assignees:
  - zhugeliang

# 是否压缩 commit
squash: true
# 删除源分支
shouldRemoveSourceBranch: true
# 需要 pipeline 成功才能合并。 注意: 无 pipeline 视为不成功
mergeWhenPipelineSucceeds: true
```

### Create Merge Request

**Commit Message** 选择使用 MR 的 `Title` 和 `Description` 是为了约束用户更加重视信息的 **规范性** 及 **完整性**。
**文档即注释**

1. `merge_commit_message` 的内容为 `merge request` 的标题 **Title**。
2. `squash_commit_message` 的内容为 `merge request` 的描述内容 **Description**。

![commit-message](./docs/mergebot/commit-message.png)

**Assignee**


2. 将 **Assignee** 选择人为 **mergebot**

![assignee-bot](./docs/mergebot/assignee-bot.png)


### Code Review

1. 任意用户都可以在 MR 下进行评论
2. 只有在 **目标分支** 中 `.mergebot.yml` 分支中的 CodeReivewer 中的任意候选人评论 `@gitlab lgtm` 时， 机器人才会进行有效计数。


### Transfer Assignee or Merge

当 `.mergebot.yml` 中的 `assignees` 控制该行为。

1. 当 `assignees` 有效时， `mergebot` 将会随机选择一位候选人， 并设置为 `Assignee` 对象。 并有改用户进行后续合并管理。
2. 当 `assignees` 无效时（无候选人）， `mergebot` 将直接进行 **合并操作**。
