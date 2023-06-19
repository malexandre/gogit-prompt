# Git prompt build with Go

I had a git prompt built with Bash, but I wanted to use a language with easier string and variables manipulations. It was important to still have a prompt with great performance, so I used Go.

```bash
main R[-0|+1]  # If you're on the [default] main branch, only diff with your current remote, origin/[default]
feature* R[-1|+3] M[-2|+3]  # On a different branch with its own remote branch, diff with both the remote R and default main M
new-feature M[-2|+3]  # On a different branch without its own remote branch, diff with only the default main M
main*  # When the repo is only local, no diff
```

## Installation

Clone this project and then run `go build`. Make sure the binary is in your path, and then call it in your custom bash prompt with `$(gogit-prompt)`.
