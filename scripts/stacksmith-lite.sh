#!/bin/bash

# stacksmith - Ultralight Artisan Git Stacking Tool (Forgive the corny jokes - have to stay on brand 🧑🏾‍🏭)

command=$1
shift

# Colors for fun output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

smith_echo() {
  echo -e "${GREEN}stacksmith${NC} 🔨 $1"
}

case $command in
  stack)
    new_branch=$1
    parent_branch=$2
    git checkout -b $new_branch $parent_branch
    smith_echo "🪵 Forged new branch $new_branch atop $parent_branch. 🌲"
    ;;

  sync)
    smith_echo "🧽 Polishing your branch stack... 🪞"
    branches=("$@")

    for (( idx=1; idx<${#branches[@]}; idx++ ))
    do
      child=${branches[$idx]}
      parent=${branches[$idx-1]}

      smith_echo "🔄 Rebasing $child onto $parent..."
      git checkout $child
      git fetch
      git rebase $parent
      git push --force-with-lease
      smith_echo "🚀 Pushed $child upstream."
    done
    ;;

  fix-pr)
    branch=$1
    target=$2

    smith_echo "🔧 Reworking $branch onto $target... 🪄"
    git checkout $branch
    git fetch
    git rebase origin/$target
    git push --force-with-lease

    smith_echo "📢 Don't forget to retarget the PR for $branch to $target in Azure DevOps, GitHub, or whatever you are using!"
    ;;

  push)
    current_branch=$(git rev-parse --abbrev-ref HEAD)

    if git rev-parse --abbrev-ref --symbolic-full-name @{u} >/dev/null 2>&1; then
      smith_echo "⬆️  Lifting $current_branch to remote forge..."
      git push --force-with-lease
    else
      smith_echo "🆕 First lift for $current_branch — setting upstream..."
      git push --set-upstream origin $current_branch --force-with-lease
    fi
    ;;

  graph)
    smith_echo "🌳 Behold your branching masterpiece:"
    git log --graph --oneline --decorate --all
    ;;

  help|*)
    echo -e "${GREEN}stacksmith${NC} - Ultralight Artisan Git Stacking Tool"
    echo ""
    echo "Usage: bash stacksmith.sh <command> [args]"
    echo ""
    echo "Commands:"
    echo "  stack <new-branch> <parent-branch>     🪵 Forge a new stacked branch"
    echo "  sync <branch1> <branch2> ...           🧽 Rebase your entire branch stack"
    echo "  fix-pr <branch> <new-target>           🔧 Rebase + remind to retarget PR"
    echo "  push                                   ⬆️  Safe push with upstream handling"
    echo "  graph                                  🌳 Visualize branch structure"
    echo "  help                                   ℹ️  Show this help menu"
    ;;
esac
