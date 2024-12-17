package junit

import (
	"testing"

	"github.com/hashicorp/terraform/internal/moduletest"
)

func Test_JUnitXMLTestReport(t *testing.T) {
	cases := map[string]struct {
		Suite     *moduletest.Suite
		XmlString string
	}{
		"no tests": {
			XmlString: "<?xml version=\"1.0\" encoding=\"UTF-8\"?><testsuites></testsuites>",
			Suite:     &moduletest.Suite{},
		},
		"one passing test": {
			XmlString: `<?xml version="1.0" encoding="UTF-8"?><testsuites>
  <testsuite name="test_name.tftest.hcl" tests="1" skipped="0" failures="0" errors="0">
    <testcase name="test_one" classname="test_name.tftest.hcl"></testcase>
  </testsuite>
</testsuites>`,
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_name.tftest.hcl": {
						Name:   "test_name.tftest.hcl",
						Status: moduletest.Skip,
						Runs: []*moduletest.Run{
							{
								Name:   "test_one",
								Status: moduletest.Pass,
							},
						},
					},
				},
			},
		},
		"one skipped test": {
			XmlString: `<?xml version="1.0" encoding="UTF-8"?><testsuites>
  <testsuite name="test_name.tftest.hcl" tests="1" skipped="1" failures="0" errors="0">
    <testcase name="test_one" classname="test_name.tftest.hcl">
      <skipped></skipped>
    </testcase>
  </testsuite>
</testsuites>`,
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_name.tftest.hcl": {
						Name:   "test_name.tftest.hcl",
						Status: moduletest.Skip,
						Runs: []*moduletest.Run{
							{
								Name:   "test_one",
								Status: moduletest.Skip,
							},
						},
					},
				},
			},
		},
		"one failed test": {
			XmlString: `<?xml version="1.0" encoding="UTF-8"?><testsuites>
  <testsuite name="test_name.tftest.hcl" tests="1" skipped="0" failures="1" errors="0">
    <testcase name="test_one" classname="test_name.tftest.hcl">
      <failure message="Test run failed"></failure>
    </testcase>
  </testsuite>
</testsuites>`,
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_name.tftest.hcl": {
						Name:   "test_name.tftest.hcl",
						Status: moduletest.Skip,
						Runs: []*moduletest.Run{
							{
								Name:   "test_one",
								Status: moduletest.Fail,
							},
						},
					},
				},
			},
		},
		"three tests, each different status": {
			XmlString: `<?xml version="1.0" encoding="UTF-8"?><testsuites>
  <testsuite name="test_name.tftest.hcl" tests="3" skipped="1" failures="1" errors="0">
    <testcase name="test_one" classname="test_name.tftest.hcl"></testcase>
    <testcase name="test_two" classname="test_name.tftest.hcl">
      <skipped></skipped>
    </testcase>
    <testcase name="test_three" classname="test_name.tftest.hcl">
      <failure message="Test run failed"></failure>
    </testcase>
  </testsuite>
</testsuites>`,
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_name.tftest.hcl": {
						Name:   "test_name.tftest.hcl",
						Status: moduletest.Skip,
						Runs: []*moduletest.Run{
							{
								Name:   "test_one",
								Status: moduletest.Pass,
							},
							{
								Name:   "test_two",
								Status: moduletest.Skip,
							},
							{
								Name:   "test_three",
								Status: moduletest.Fail,
							},
						},
					},
				},
			},
		},
		"multiple test files with various tests": {
			XmlString: `<?xml version="1.0" encoding="UTF-8"?><testsuites>
  <testsuite name="test_file_one.tftest.hcl" tests="1" skipped="0" failures="0" errors="0">
    <testcase name="test_one" classname="test_file_one.tftest.hcl"></testcase>
  </testsuite>
  <testsuite name="test_file_two.tftest.hcl" tests="2" skipped="1" failures="1" errors="0">
    <testcase name="test_two" classname="test_file_two.tftest.hcl">
      <skipped></skipped>
    </testcase>
    <testcase name="test_three" classname="test_file_two.tftest.hcl">
      <failure message="Test run failed"></failure>
    </testcase>
  </testsuite>
</testsuites>`,
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_file_one.tftest.hcl": {
						Name:   "test_file_one.tftest.hcl",
						Status: moduletest.Skip,
						Runs: []*moduletest.Run{
							{
								Name:   "test_one",
								Status: moduletest.Pass,
							},
						},
					},
					"test_file_two.tftest.hcl": {
						Name:   "test_file_two.tftest.hcl",
						Status: moduletest.Skip,
						Runs: []*moduletest.Run{
							{
								Name:   "test_two",
								Status: moduletest.Skip,
							},
							{
								Name:   "test_three",
								Status: moduletest.Fail,
							},
						},
					},
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			b, _ := JUnitXMLTestReport(tc.Suite)
			if string(b) != tc.XmlString {
				t.Fatalf("wanted XML:\n%s\n got XML:\n%s\n", tc.XmlString, string(b))
			}
		})
	}
}

func Test_suiteFilesAsSortedList(t *testing.T) {
	cases := map[string]struct {
		Suite         *moduletest.Suite
		ExpectedNames map[int]string
	}{
		"no test files": {
			Suite: &moduletest.Suite{},
		},
		"3 test files ordered in map": {
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_file_1.tftest.hcl": {
						Name:   "test_file_1.tftest.hcl",
						Status: moduletest.Skip,
						Runs:   []*moduletest.Run{},
					},
					"test_file_2.tftest.hcl": {
						Name:   "test_file_2.tftest.hcl",
						Status: moduletest.Skip,
						Runs:   []*moduletest.Run{},
					},
					"test_file_3.tftest.hcl": {
						Name:   "test_file_3.tftest.hcl",
						Status: moduletest.Skip,
						Runs:   []*moduletest.Run{},
					},
				},
			},
			ExpectedNames: map[int]string{
				0: "test_file_1.tftest.hcl",
				1: "test_file_2.tftest.hcl",
				2: "test_file_3.tftest.hcl",
			},
		},
		"3 test files unordered in map": {
			Suite: &moduletest.Suite{
				Status: moduletest.Skip,
				Files: map[string]*moduletest.File{
					"test_file_3.tftest.hcl": {
						Name:   "test_file_3.tftest.hcl",
						Status: moduletest.Skip,
						Runs:   []*moduletest.Run{},
					},
					"test_file_1.tftest.hcl": {
						Name:   "test_file_1.tftest.hcl",
						Status: moduletest.Skip,
						Runs:   []*moduletest.Run{},
					},
					"test_file_2.tftest.hcl": {
						Name:   "test_file_2.tftest.hcl",
						Status: moduletest.Skip,
						Runs:   []*moduletest.Run{},
					},
				},
			},
			ExpectedNames: map[int]string{
				0: "test_file_1.tftest.hcl",
				1: "test_file_2.tftest.hcl",
				2: "test_file_3.tftest.hcl",
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			list := suiteFilesAsSortedList(tc.Suite.Files)

			if len(tc.ExpectedNames) != len(tc.Suite.Files) {
				t.Fatalf("expected there to be %d items, got %d", len(tc.ExpectedNames), len(tc.Suite.Files))
			}

			if len(tc.ExpectedNames) == 0 {
				return
			}

			for k, v := range tc.ExpectedNames {
				if list[k].Name != v {
					t.Fatalf("expected element %d in sorted list to be named %s, got %s", k, v, list[k].Name)
				}
			}
		})
	}
}
