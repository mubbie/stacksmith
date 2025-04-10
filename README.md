# Stacksmith

> Ultralight Artisan Git Stacking Tool (Forgive the corny jokes - have to stay on brand 🧑🏾‍🏭)

A tiny Bash-powered tool for developers managing stacked pull requests using vanilla Git. 🌳

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

## 🤝 Contribution

Contributions are welcome! ✨
1. Fork the project
2. Create your feature branch (git checkout -b feat/amazing-feature)
3. Commit your changes (git commit -m 'feat: add amazing feature')
4. Push to the branch (git push)
5. Open a pull request
