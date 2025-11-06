package phash

import "testing"

func TestFnvSignature(t *testing.T) {
	label := map[string]string{
		"name": "jack",
		"age":  "18",
	}

	x := LabelsToSignature(label).String()
	t.Log(x)
	t.Log(FingerprintFromString(x))

	d := LabelsToSignature(label).Format()
	t.Log(d)
	t.Log(FingerprintFromFormat(d))
}
