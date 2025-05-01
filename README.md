# Contribution Guide

### Git conventions
Branch Naming
| Type       | Purpose |
|------------|---------|
| `feature/` | New features or enhancements |
| `bugfix/`  | Non-critical bug fixes |
| `hotfix/`  | Critical fixes that go directly to production |
| `chore/`   | Routine tasks like dependency updates |
| `refactor/`| Code restructuring without changing behavior |
| `test/`    | Adding or updating tests |
| `docs/`    | Documentation changes only |

e.g. `feature/user-signup` `bugfix/fix-null-error`
if Jira card exists `feature/PROJ-123-login-flow`

Commit Message
| Type       | Purpose |
|------------|---------|
| `feat`     | A new feature |
| `fix`      | A bug fix |
| `docs`     | Documentation only changes |
| `style`    | Formatting, missing semicolons, etc (no code changes) |
| `refactor` | Code change that doesn't fix a bug or add a feature |
| `perf`     | Performance improvement |
| `test`     | Adding or modifying tests |
| `chore`    | Routine tasks like dependency updates |
| `ci`       | CI/CD pipeline changes |
| `build`    | Build system or external dependencies changes |

```
feat(auth): add login endpoint
fix(user): prevent crash on null profile
docs(readme): add setup instructions
style(linter): fix ESLint errors
refactor(cart): simplify discount logic
test(api): add test for payment gateway
chore(deps): update lodash to v4.17.21
```
