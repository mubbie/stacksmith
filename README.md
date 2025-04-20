# Stacksmith

> Ultralight Artisan Git Stacking Tool (Forgive the corny jokes — staying on brand 🧑🏾‍🏭)
> *Forgive the corny jokes — staying on brand 🧑🏾‍🏭*

Stacksmith is your terminal blacksmithing forge for managing stacked pull branches and pull requests using **vanilla Git** 🌳

This repo contains two versions:

| Version             | Description                                         |
|----------------------|-----------------------------------------------------|
| `stacksmith-lite.sh` | 🪶 Lightweight Bash script for fast Git stacking    |
| `stacksmith`         | ⚡ Upcoming Go-powered CLI with rich UI (coming soon) |

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

### What Are Stacked PRs? 🚂

Stacked PRs let you break work into small, focused branches — each building on top of the last.

```text
main <- PR1 <- PR2 <- PR3 <- PR4 ...
```

Each PR targets the previous one, reviewers see small diffs, and you keep moving fast.

BUT managing these stacks manually with plain Git is tedious (See: [Stacked branches with vanilla Git](https://www.codetinkerer.com/2023/10/01/stacked-branches-with-vanilla-git.html), [Stacked branches with vanilla Git - Reddit Thread](https://www.reddit.com/r/programming/comments/16yqfef/stacked_branches_with_vanilla_git/)):

- Rebasing every branch on top of the latest
- Force pushing without messing things up
- Retargeting PRs

That's where `stacksmith` comes in.

---

## Stacksmith Lite 🪶

`stacksmith-lite.sh` is a zero-installation, dead-simple Bash script for managing stacked branches using **vanilla Git.**

It works anywhere Git works:

- ✅ Local dev
- ✅ CI environments
- ✅ Remote VMs
- ✅ No plugins, no wrappers, no setup

### Install Stacksmith Lite 🚀

```bash
curl -sL https://raw.githubusercontent.com/mubbie/stacksmith/main/scripts/stacksmith-lite.sh -o stacksmith
chmod +x stacksmith
sudo mv stacksmith /usr/local/bin/
```

Or just alias it:

```bash
alias stacksmith='bash /path/to/stacksmith-lite.sh'
```

If you run into trouble adding `stacksmith` to your path, [here's](https://specifications.freedesktop.org/basedir-spec/latest/) an excellent and helpful article recommended by my friend [Osaro](https://github.com/osaroadade) 🙂

### Usage ⚙️

#### 🪵 Create a new stacked branch

```bash
stacksmith stack <new-branch> <parent-branch>
```

#### 🧽 Rebase and sync your stack

```bash
stacksmith sync <branch1> <branch2> <branch3> ...
```

#### 🔧 Rebase a branch after parent PR merges

```bash
stacksmith fix-pr <branch> <new-target>
```

#### ⬆️ Push current branch safely

```bash
stacksmith push
```

#### 🌳 Visualize your branch stack

```bash
stacksmith graph
```

---

## Managing PRs with Stacksmith 📂

> Stacksmith helps you manage your local branches beautifully. But your PRs will still need to be created, managed, and merged manually on your Git hosting platform (Azure DevOps, GitHub, GitLab, Bitbucket, etc).

### PR Lifecycle with Stacksmith

- Create your stacked branches locally with `stacksmith stack`
- Push them with `stacksmith push`
- Open PRs in your Git platform (targeting their parent branches, ex: ex: PR2 targets PR1, PR3 targets PR2, etc.)
- Merge PRs bottom-up (base first, then next, then next)
- After each PR merge:
  - Use `stacksmith fix-pr` to rebase the next branch onto the new target (usually `main`)
  - Retarget the PR in your Git platform to point to `main`
  - Push again with `stacksmith push`

### Pro Tip

Use `stacksmith sync` to quickly rebase and update a full stack when many PRs have merged.

- Stacksmith = Local branch management magic
- Your Git platform = PR creation, review, merging
- Together = Dev happiness 🌟

### What Stacksmith Doesn't Do 🙅

- ❌ Create PRs for you (use your Git platform)
- ❌ Auto-retarget PRs (you do that manually)
- ❌ Auto-detect your stack (you pass branch names explicitly)

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

## Comming Soon: Stacksmith (Go Edition) ⚡

We’re rebuilding Stacksmith in Go for a more powerful and visual CLI experience:
- 🌲 Rich, colorized DAG views of your stack
- 🧑🏾‍🏭 Interactive flows powered by [gum](https://github.com/charmbracelet/gum)
- 💻 Optional full-screen TUI with [bubbletea](https://github.com/charmbracelet/bubbletea)
- 🧪 Diff previews, merge awareness, branch introspection
- 🧼 Safe, guided stack syncing and push flows

It’ll be fast, easy to install, and compatible with your current `stacksmith-lite` setup. Stay tuned!

---

## 🤝 Contribution

Contributions are welcome! ✨
1. Fork the project
2. Create your feature branch (git checkout -b feat/amazing-feature)
3. Commit your changes (git commit -m 'feat: add amazing feature')
4. Push to the branch (git push)
5. Open a pull request

Add commands, fix bugs, clean up UI, or just drop a pun. All artisan hands on deck.

---

## 📢 Feedback

Got ideas, bugs, or thoughts? Love a bad artisan pun? Open an issue or reach out!

Your feedback makes this tool better (and funnier). 😎

---

Crafted with love (and corny jokes), by artisans of the stack (mostly GPT-4o). 🧑🏾‍🏭✨
