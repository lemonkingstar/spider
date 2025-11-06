package phash

import (
	"fmt"
	"sort"
	"strconv"
)

// SeparatorByte is a byte that cannot occur in valid UTF-8 sequences and is
// used to separate label names, label values, and other strings from each other
// when calculating their combined hash value (aka signature aka fingerprint).
const SeparatorByte byte = 255

// Fingerprint provides a hash-capable representation of a Metric.
// For our purposes, FNV-1A 64-bit is used.
type Fingerprint uint64

var (
	// cache the signature of an empty label set.
	emptyLabelSignature = hashNew()
)

// LabelsToSignature returns a quasi-unique signature (i.e., fingerprint) for a
// given label set. (Collisions are possible but unlikely if the number of label
// sets the function is applied to is small.)
func LabelsToSignature(labels map[string]string) Fingerprint {
	if len(labels) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	labelNames := make([]string, 0, len(labels))
	for labelName := range labels {
		labelNames = append(labelNames, labelName)
	}
	sort.Strings(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames {
		sum = hashAdd(sum, labelName)
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, labels[labelName])
		sum = hashAddByte(sum, SeparatorByte)
	}
	return Fingerprint(sum)
}

func (f Fingerprint) String() string {
	// 16进制16位输出 不够补零
	return fmt.Sprintf("%016x", uint64(f))
}

func (f Fingerprint) Format() string {
	// 十进制格式化
	return fmt.Sprintf("%d", uint64(f))
}

// FingerprintFromString transforms a string representation into a Fingerprint.
func FingerprintFromString(s string) (uint64, error) {
	num, err := strconv.ParseUint(s, 16, 64)
	return num, err
}

// FingerprintFromFormat transforms a string representation into a Fingerprint.
func FingerprintFromFormat(s string) (uint64, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return num, err
}
