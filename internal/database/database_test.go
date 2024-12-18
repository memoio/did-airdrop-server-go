package database

import (
	"os"
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
)

func TestGetNumber(t *testing.T) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	db, err := CreateDB(klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	num, err := db.GetNumber()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(num)
}

func TestAddNumber(t *testing.T) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	db, err := CreateDB(klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddNumber("did:memo:947e38821cec0d483922bf082958caa38c9c8900cdd9184a159ea07a5e18b9ac", 100001)
	if err != nil {
		t.Fatal(err)
	}
}
