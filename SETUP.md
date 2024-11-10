## Setting up Bazel

1. Create a MODULE.bazel with the following content
```python
bazel_dep(name = "rules_go", version = "0.50.1")
bazel_dep(name = "gazelle", version = "0.39.1")

go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//:go.mod")
use_repo(go_deps, "com_github_stretchr_testify")
```

2. Add the following to your BUILD.bazel file
```python
load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_library", "go_test")

gazelle(name = "gazelle")

# gazelle:prefix github.com/minhajuddin/injector3
```

3. Run
```bash
bazel run //:gazelle
```

4. Run go mod tidy
```bash
go mod tidy
```

5. Run bazel mod tidy
```bash
bazel mod tidy
```

6. Run the tests
```bash
bazel test //... --test_output=all --flaky_test_attempts=3
```
