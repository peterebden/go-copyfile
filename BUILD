cgo_library(
    name = "copyfile",
    srcs = [
        "cow_linux.go" if CONFIG.OS == "linux" else "cow_darwin.go" if CONFIG.OS == "darwin" else "cow_other.go",
    ],
    go_srcs = ["copier.go"],
)

cgo_test(
    name = "copy_test",
    srcs = ["copy_test.go"],
    data = ["test_data"],
    deps = [
        ":copyfile",
        ":testify",
    ],
)

go_library(
    name = "copyfile_generic",
    srcs = ["copier.go", "cow_other.go"],
)

go_test(
    name = "copy_generic_test",
    srcs = ["copy_test.go"],
    data = ["test_data"],
    deps = [
        ":copyfile",
        ":testify",
    ],
)

go_get(
    name = "testify",
    get = "github.com/stretchr/testify/...",
    revision = "f390dcf405f7b83c997eac1b06768bb9f44dec18",
    install = [
        "assert",
        "require",
        "vendor/github.com/davecgh/go-spew/spew",
        "vendor/github.com/pmezard/go-difflib/difflib",
    ],
    test_only = True,
)
