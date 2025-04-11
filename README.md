# Stacksmith

> Ultralight Artisan Git Stacking Tool (Forgive the corny jokes — staying on brand 🧑🏾‍🏭)

A tiny Bash-powered tool for developers managing stacked pull requests using vanilla Git. 🌳

---

## Why Stacksmith? 🤔

Imagine this: you're building a big feature. It's going to touch a lot of files and introduce a lot of changes.
With traditional Git workflows, your options are usually:

- 🫠 Put it all in one huge branch → easy for you, painful for your reviewers.
- ⏳ Break it into many small PRs → good for reviewers, but you end up stuck waiting for each PR to merge before you can build on the next one.

Both kinda suck.

```text
Option 1 → One giant PR 😱

Option 2 → Many PRs but blocked 😩

Stacksmith → Many PRs. Keep shipping 🚀
```

**Enter Stacked PRs: The Best of Both Worlds 🚂**

Stacked PRs let you:

1. Break your work into smaller, easy-to-review pieces.
2. Stack branches on top of each other like this:

```text
main <- PR1 <- PR2 <- PR3 <- PR4 ...
```

3. Reviewers get focused diffs. You keep moving fast.

BUT managing these stacks manually with plain Git is tedious (See existing recommendations on how to approach it: [Stacked branches with vanilla Git](https://www.codetinkerer.com/2023/10/01/stacked-branches-with-vanilla-git.html), [Stacked branches with vanilla Git - Reddit Thread](https://www.reddit.com/r/programming/comments/16yqfef/stacked_branches_with_vanilla_git/)):

- Rebasing every branch on top of the latest
- Force pushing without messing things up
- Retargeting PRs

**Aren't there tools for this already? 🤓**

Yes! There are great tools out there like:

- [Graphite](https://graphite.dev/) — excellent, powerful, but heavily tied to GitHub and its own ecosystem.
- [GitButler](https://gitbutler.com/) — super promising, but still evolving and platform-dependent.
- GitHub's own support for stacked PRs — only partial and GitHub-specific.

But sometimes...

- You want to use plain ol' Git (Turned 20 recently, see interesting interview with the creator of Git and Linux, Linus Torvalds: [Two decades of Git: A conversation with creator Linus Torvalds](https://www.youtube.com/watch?v=sCr_gb8rdEI)
- You want a tool that's portable, bash-native, no setup, no login, no install headache.
- You want something that's easy to teach, easy to adopt, works anywhere.

That's what `stacksmith` is for 🧑🏾‍🏭

A tiny, dead-simple, artisan-crafted bash tool for anyone who wants the superpower of stacked PRs without the weight of extra platforms or tools.

---

## Managing PRs with Stacksmith 📂

> Stacksmith helps you manage your local branches beautifully. But your PRs will still need to be created, managed, and merged manually on your Git hosting platform (Azure DevOps, GitHub, GitLab, Bitbucket, etc).

### PR Lifecycle with Stacksmith:

- Create your stacked branches locally with `stacksmith stack`
- Push them with `stacksmith push`
- Open PRs in your Git platform (targeting their parent branches, ex: ex: PR2 targets PR1, PR3 targets PR2, etc.)
- Merge PRs bottom-up (base first, then next, then next)
- After each PR merge:
  - Use `stacksmith fix-pr` to rebase the next branch onto the new target (usually `main`)
  - Retarget the PR in your Git platform to point to `main`
  - Push again with `stacksmith push`

### Pro Tip:

Use `stacksmith sync` to quickly rebase and update a full stack when many PRs have merged.

- Stacksmith = Local branch management magic
- Your Git platform = PR creation, review, merging
- Together = Dev happiness 🌟

### What Stacksmith Doesn't Do 🙅

- Create PRs for you (use your Git platform)
- Auto-retarget PRs (you do that manually)
- Auto-detect your stack (you pass branch names explicitly)

Stacksmith stays simple & bashy — that's the point.

---

## 🔦 Gotchas & Pitfalls

> Some common sharp edges when working with stacked PRs (and how to avoid them):

| Situation                          | What Happens                                              | How To Handle                                                                |
| ---------------------------------- | --------------------------------------------------------- | ---------------------------------------------------------------------------- |
| PR merges out of order             | Git history gets messy; later PR shows unexpected changes | Rebase your branch onto `main` using `stacksmith fix-pr` and retarget the PR |
| Forgetting to retarget PR          | PR shows extra unrelated commits                          | Always retarget PR to `main` (or the correct parent) after parent merges     |
| Not force-pushing after rebase     | Remote branch gets out of sync with local                 | Always use `stacksmith push` (safe force-push) after rebasing                |
| Accidentally rebasing wrong parent | Changes vanish or conflict                                | Double-check the branch order when using `stacksmith sync`                   |

**Final Rule of Thumb:**

- Merge PRs from the bottom up
- Rebase child branches immediately after parent merges
- Retarget PRs accordingly
- Push your changes
- Clean stack = Happy reviewers + Happy you 🌱

---

## 🚀 Installation

1. Download `stacksmith.sh`
2. Make it executable:
```bash
chmod +x stacksmith.sh
```
3. (Optional) Add to PATH:
```bash
mv stacksmith.sh ~/bin/stacksmith
```
Or alias it:
```
alias stacksmith='bash /path/to/stacksmith.sh'
```

---

## ⚙️ Usage

### 🪵 Create a new stacked branch

```bash
stacksmith stack <new-branch> <parent-branch>
```

Example:

```bash
stacksmith stack feature/part-1 feature/base
```

### 🧽 Rebase and sync your stack

```bash
stacksmith sync <branch1> <branch2> <branch3> ...
```

Example:

```bash
stacksmith sync feature/base feature/part-1 feature/part-2
```

### 🔧 Rebase a branch after parent PR merges

```bash
stacksmith fix-pr <branch> <new-target>
```

Example:

```bash
stacksmith fix-pr feature/part-1 main
```

### ⬆️ Push current branch safely

```bash
stacksmith push
```

Handles first-time push & force push safely.

### 🌳 Visualize your branch stack

```bash
stacksmith graph
```

---

## 🤝 Contribution

Contributions are welcome! ✨
1. Fork the project
2. Create your feature branch (git checkout -b feat/amazing-feature)
3. Commit your changes (git commit -m 'feat: add amazing feature')
4. Push to the branch (git push)
5. Open a pull request

---

## 📢 Feedback

Got ideas, bugs, or thoughts? Love a bad artisan pun? Open an issue or reach out!

Your feedback makes this tool better (and funnier). 😎

---

Crafted with love (and corny jokes), by artisans of the stack (mostly GPT-4o). 🧑🏾‍🏭✨
