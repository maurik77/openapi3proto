package openapi2proto_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/NYTimes/openapi2proto"
	"github.com/NYTimes/openapi2proto/compiler"
	"github.com/NYTimes/openapi2proto/protobuf"
	"github.com/pmezard/go-difflib/difflib"
)

type genProtoTestCase struct {
	options                 bool
	fixturePath             string
	wantProto               string
	remoteFiles             []string
	wrapPrimitives          bool
	skipDeprecatedRpcs      bool
	addAutogeneratedComment bool
}

func testGenProto(t *testing.T, tests ...genProtoTestCase) {
	t.Helper()
	origin, _ := os.Getwd()
	for _, test := range tests {
		t.Run(test.fixturePath, func(t *testing.T) {
			for _, remoteFile := range test.remoteFiles {
				res, err := http.Get(remoteFile)
				if err != nil || res.StatusCode != http.StatusOK {
					t.Skip(`Remote file ` + remoteFile + ` is not available`)
				}
			}

			var generated bytes.Buffer
			var compilerOptions []compiler.Option
			var encoderOptions []protobuf.Option
			if test.options {
				compilerOptions = append(compilerOptions, compiler.WithAnnotation(true))
			}
			if test.wrapPrimitives {
				compilerOptions = append(compilerOptions, compiler.WithWrapPrimitives(true))
			}
			if test.skipDeprecatedRpcs {
				compilerOptions = append(compilerOptions, compiler.WithSkipDeprecatedRpcs(true))
			}
			if test.addAutogeneratedComment {
				encoderOptions = append(encoderOptions, protobuf.WithAutogeneratedComment(true))
			}
			if err := openapi2proto.Transpile(&generated, test.fixturePath, openapi2proto.WithCompilerOptions(compilerOptions...), openapi2proto.WithEncoderOptions(encoderOptions...)); err != nil {
				t.Errorf(`failed to transpile: %s`, err)
				return
			}

			os.Chdir(origin)
			// if test.wantProto is empty, guess file name from the original
			// fixture path
			wantProtoFile := test.wantProto
			if wantProtoFile == "" {
				i := strings.LastIndexByte(test.fixturePath, '.')
				if i > -1 {
					wantProtoFile = test.fixturePath[:i] + `.proto`
				} else {
					t.Fatalf(`unable to guess proto file name from %s`, test.fixturePath)
				}
			}
			want, err := ioutil.ReadFile(wantProtoFile)
			if err != nil {
				t.Fatal("unable to open test fixture: ", err)
			}

			if string(want) != generated.String() {
				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(string(want)),
					B:        difflib.SplitLines(generated.String()),
					FromFile: wantProtoFile,
					ToFile:   "Generated",
					Context:  3,
				}
				text, _ := difflib.GetUnifiedDiffString(diff)
				t.Errorf("testYaml (%s) differences:\n%s",
					test.fixturePath, text)
			}
		})
	}
}

func TestNetwork(t *testing.T) {
	testGenProto(t, genProtoTestCase{
		fixturePath: "fixtures/petstore/swagger.yaml",
		remoteFiles: []string{
			"https://raw.githubusercontent.com/NYTimes/openapi2proto/master/fixtures/petstore/Pet.yaml",
		},
	})
}

func TestGenerateProto(t *testing.T) {
	tests := []genProtoTestCase{
		// {
		// 	fixturePath: "fixtures/cats.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/catsanddogs.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/semantic_api.json",
		// },
		// {
		// 	fixturePath: "fixtures/semantic_api.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/most_popular.json",
		// },
		// {
		// 	fixturePath: "fixtures/spec.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/spec.json",
		// },
		// {
		// 	options:     true,
		// 	fixturePath: "fixtures/semantic_api.json",
		// 	wantProto:   "fixtures/semantic_api-options.proto",
		// },
		// {
		// 	options:     true,
		// 	fixturePath: "fixtures/most_popular.json",
		// 	wantProto:   "fixtures/most_popular-options.proto",
		// },
		// {
		// 	options:     true,
		// 	fixturePath: "fixtures/spec.yaml",
		// 	wantProto:   "fixtures/spec-options.proto",
		// },
		// {
		// 	options:     true,
		// 	fixturePath: "fixtures/spec.json",
		// 	wantProto:   "fixtures/spec-options.proto",
		// },

		// {
		// 	fixturePath: "fixtures/includes_query.json",
		// },
		// {
		// 	fixturePath: "fixtures/lowercase_def.json",
		// },
		// {
		// 	fixturePath: "fixtures/missing_type.json",
		// },
		// /*
		// 	{
		// 		fixturePath: "fixtures/kubernetes.json",
		// 	},
		// */
		// {
		// 	fixturePath: "fixtures/accountv1-0.json",
		// },
		// // {
		// // 	fixturePath: "fixtures/refs.json",
		// // },
		// {
		// 	fixturePath: "fixtures/refs.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/integers.yaml",
		// },
		// {
		// 	wrapPrimitives: true,
		// 	fixturePath:    "fixtures/integers_required.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/global_options.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/naming_conversion.yaml",
		// },
		// {
		// 	options:     true,
		// 	fixturePath: "fixtures/custom_options.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/string_proto_tag.yaml",
		// },
		// {
		// 	skipDeprecatedRpcs: true,
		// 	fixturePath:        "fixtures/skip_deprecated_rpcs.yaml",
		// },
		// {
		// 	addAutogeneratedComment: true,
		// 	fixturePath:             "fixtures/add_autogenerated_comment.yaml",
		// },
		// {
		// 	fixturePath: "fixtures/global_responses.yaml",
		// },
		{
			fixturePath: "fixtures/spec_v3.yaml",
			wantProto:   "fixtures/spec.proto",
		},
	}
	testGenProto(t, tests...)
}
