go_library(
    name = "copyfile",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
)

go_test(
    name = "copy_test",
    srcs = ["copy_test.go"],
    data = ["test_data"],
    deps = [
        ":testify",
    ],
)

go_get(
    name = "spew",
    get = "github.com/davecgh/go-spew/spew",
    revision = "ecdeabc65495df2dec95d7c4a4c3e021903035e5",
    test_only = True,
)

go_get(
    name = "testify",
    get = "github.com/stretchr/testify/...",
    revision = "f390dcf405f7b83c997eac1b06768bb9f44dec18",
    test_only = True,
    deps = [":spew"],
)
